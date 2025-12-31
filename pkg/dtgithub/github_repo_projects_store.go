package dtgithub

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/datatug/datatug-core/pkg/datatug"
	"github.com/datatug/datatug-core/pkg/dtconfig"
	"github.com/datatug/datatug-core/pkg/storage/filestore"
	"github.com/google/go-github/v80/github"
	"gopkg.in/yaml.v3"
)

func NewRepoProjectsStore(client *github.Client, branch string) *GithubRepoProjectsStore {
	if branch == "" {
		branch = "main"
	}
	return &GithubRepoProjectsStore{
		client: client,
		branch: branch,
	}
}

var _ datatug.ProjectsStore = (*GithubRepoProjectsStore)(nil)

type GithubRepoProjectsStore struct {
	client *github.Client
	branch string
}

func (g GithubRepoProjectsStore) CreateNewProject(
	ctx context.Context,
	projectID, title string,
	visibility datatug.ProjectVisibility,
	report datatug.StatusReporter,
) (
	project *datatug.Project, err error,
) {
	ids := strings.Split(projectID, "/")
	repoOwner, repoName, projectDir := ids[0], ids[1], path.Join(ids[2:]...)
	_, _, _ = repoOwner, repoName, title

	pathToProjectFromRepoRoot := path.Join(projectDir)

	var entries []*github.TreeEntry

	addEntry := func(path string, content string) {
		// Check if files already exist to avoid overwriting or redundant commits
		if existing, _, _, _ := g.client.Repositories.GetContents(ctx, repoOwner, repoName, path, &github.RepositoryContentGetOptions{Ref: g.branch}); existing == nil {
			entries = append(entries, &github.TreeEntry{
				Path:    github.Ptr(path),
				Type:    github.Ptr("blob"),
				Mode:    github.Ptr("100644"),
				Content: github.Ptr(content),
			})
		}
	}

	creator := newProjectCreator(g.client, report, addEntry)
	if err = creator.CreateProject(ctx, title, pathToProjectFromRepoRoot, visibility); err != nil {
		return
	}

	if len(entries) > 0 {
		//var tree *github.Tree
		//tree, _, err = g.client.Git.CreateTree(ctx, repoOwner, repoName, *ref.Object.SHA, entries)
		//if err != nil {
		//	err = fmt.Errorf("failed to create tree: %w", err)
		//	return
		//}
		panic("not implemented")
	}

	if isCancelled(ctx) {
		return
	}

	time.Sleep(500 * time.Millisecond)

	if isCancelled(ctx) {
		return
	}

	//tui.App.QueueUpdateDraw(func() {
	//	loader := filestore.NewProjectsLoader("~/datatug")
	//	pConfig := dtconfig.ProjectRef{
	//		ID:    projectID,
	//		Title: title,
	//		Path:  projectID,
	//	}
	//	projectCtx := dtproject.NewProjectContext(tui, loader, &pConfig)
	//	GoProjectScreen(projectCtx)
	//})
	return
}

type projectCreator struct {
	client         *github.Client
	branch         string
	repoName       string
	repoOwner      string
	repo           *github.Repository
	addEntry       func(path, content string)
	report         datatug.StatusReporter
	steps          []*datatug.Step
	stepRepo       *datatug.Step
	stepRootFile   *datatug.Step
	stepRootReadme *datatug.Step
	stepProjFile   *datatug.Step
	stepClone      *datatug.Step
}

func (c *projectCreator) reportStatus() {
	c.report(c.steps)
}

func newProjectCreator(ghClient *github.Client, report datatug.StatusReporter, addEntry func(path, content string)) (creator *projectCreator) {
	creator = &projectCreator{
		client:         ghClient,
		addEntry:       addEntry,
		report:         report,
		stepRepo:       &datatug.Step{Name: "Create repository", Status: "pending"},
		stepProjFile:   &datatug.Step{Name: "Create .datatug-project.json", Status: "pending"},
		stepRootFile:   &datatug.Step{Name: "Update /.datatug.yaml", Status: "pending"},
		stepRootReadme: &datatug.Step{Name: "Update /README.md", Status: "pending"},
		stepClone:      &datatug.Step{Name: "Clone repo", Status: "pending"},
	}
	creator.steps = []*datatug.Step{
		creator.stepRepo,
		creator.stepRootFile,
	}
	report(creator.steps)
	return &projectCreator{
		steps: []*datatug.Step{},
	}
}

func (c *projectCreator) CreateProject(
	ctx context.Context,
	title, pathToProjectFromRepoRoot string,
	visibility datatug.ProjectVisibility,
) (err error) {

	if err = c.createRepo(ctx, visibility); err != nil {
		return fmt.Errorf("failed to create GitHub repository '%s/%s': %w", c.repoOwner, c.repoName, err)
	}

	if err = c.cloneRepo(); err != nil {
		return fmt.Errorf("failed to clone GitHub repository '%s/%s': %w", c.repoOwner, c.repoName, err)
	}

	var project *datatug.Project

	// We use the Git Data API to create multiple files in a single commit.
	// 1. Get the latest commit of the branch
	var ref *github.Reference
	ref, _, err = c.client.Git.GetRef(ctx, c.repoOwner, c.repoName, "heads/"+c.branch)
	_ = ref
	if err != nil {
		// If the repository is empty, we need to create the first commit
		var gErr *github.ErrorResponse
		if errors.As(err, &gErr) && (gErr.Response.StatusCode == 404 || gErr.Response.StatusCode == 409) {
			// Create initial README.md to initialize the repository
			_, _, err = c.client.Repositories.CreateFile(ctx, c.repoOwner, c.repoName, "README.md", &github.RepositoryContentFileOptions{
				Message: github.Ptr("feat: initial commit"),
				Content: []byte("# " + c.repoName + "\n\nDataTug project repository."),
				Branch:  github.Ptr(c.branch),
			})
			if err != nil {
				err = fmt.Errorf("failed to initialize repository: %w", err)
				return
			}
			// Retry getting the ref
			ref, _, err = c.client.Git.GetRef(ctx, c.repoOwner, c.repoName, "heads/"+c.branch)
			_ = ref
		}
		if err != nil {
			err = fmt.Errorf("failed to get branch ref: %w", err)
			return
		}
	}

	if err = c.addProjectToDataTugConfig(pathToProjectFromRepoRoot, title); err != nil {
		return fmt.Errorf("failed to add project to DataTug config: %w", err)
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var errs []error

	execute := func(f ...func() error) {
		wg.Add(len(f))
		for _, fn := range f {
			go func() {
				if fErr := fn(); fErr != nil {
					mutex.Lock()
					errs = append(errs, fErr)
					mutex.Unlock()
				}
			}()
		}
	}
	execute( // TODO: reuse parallel runner or document why not?
		func() error {
			return c.addProjectToRootRepoFile(ctx, pathToProjectFromRepoRoot)
		},
		func() error {
			return c.createProjectReadme(pathToProjectFromRepoRoot)
		},
		func() error {
			return c.createProjectSummaryFile(pathToProjectFromRepoRoot, project)
		},
	)

	wg.Wait()

	if err = c.addDatatugSectionToRootReadmeFile(ctx, pathToProjectFromRepoRoot); err != nil {
		err = fmt.Errorf("failed to add DataTug section to repo's root README.md: %w", err)
		return err
	}

	return
}

func (c *projectCreator) addProjectToRootRepoFile(_ context.Context, projPath string) error {
	c.stepRootFile.Status = "updating..."
	c.reportStatus()
	if projPath == "" {
		projPath = "."
	}
	var repoRootFile datatug.RepoRootFile
	repoRootFile.Projects = append(repoRootFile.Projects, projPath)
	content, err := yaml.Marshal(repoRootFile)
	if err != nil {
		c.stepRootFile.Status = "[red]error: " + err.Error()
		c.reportStatus()
		return fmt.Errorf("failed to marshal repoRootFile: %w", err)
	}
	filePath := path.Join(projPath, filestore.RepoRootDataTugFileName)
	c.addEntry(filePath, string(content))
	c.stepRootFile.Status = "[green]updated[-]."
	c.reportStatus()
	return nil
}

func (c *projectCreator) createRepo(ctx context.Context, visibility datatug.ProjectVisibility) (err error) {
	c.repo, _, err = c.client.Repositories.Get(context.Background(), c.repoOwner, c.repoName)
	if err != nil {
		// Create repository
		c.repo = &github.Repository{
			Name:    github.Ptr(c.repoName),
			Private: github.Ptr(visibility == datatug.PrivateProject),
		}

		c.repo, _, err = c.client.Repositories.Create(ctx, "", c.repo)
	}
	return err
}

func (c *projectCreator) addDatatugSectionToRootReadmeFile(ctx context.Context, projPath string) error {
	c.stepRootReadme.Status = "updating..."
	c.reportStatus()

	// 2. Add 'DataTug' section to root README.md
	rootReadme, _, err := c.client.Repositories.GetReadme(ctx, c.repoOwner, c.repoName, &github.RepositoryContentGetOptions{Ref: c.branch})

	var dataTugSectionTitleRegex = regexp.MustCompile(`\n##\s*DataTug`)

	getDataTugSectionForReadmeMD := func() string {
		const dataTugSectionTitleText = "DataTug - [github.com/datatug/datatug](https://github.com/datatug/datatug)"
		appLink := fmt.Sprintf("[DataTug.app](https://datatug.app/#project=github.com/%s/%s/%s)", c.repoOwner, c.repoName, projPath)
		msg := fmt.Sprintf("The [/datatug](./datatug) project can be opened and edited in %s.", appLink)
		return fmt.Sprintf("\n\n## DataTug - %s\n\n%s\n\n", dataTugSectionTitleText, msg)
	}

	if err == nil {
		content, _ := rootReadme.GetContent()
		if !dataTugSectionTitleRegex.Match([]byte(content)) {
			newContent := content + getDataTugSectionForReadmeMD()
			_, _, err = c.client.Repositories.UpdateFile(ctx, c.repoOwner, c.repoName, rootReadme.GetPath(), &github.RepositoryContentFileOptions{
				Message: github.Ptr("chore: adds ##DataTug section to /README.md"),
				Content: []byte(newContent),
				SHA:     rootReadme.SHA,
				Branch:  github.Ptr(c.branch),
			})
			if err != nil {
				return fmt.Errorf("failed to update /README.md: %w", err)
			}
		}
	} else {
		newContent := "# " + c.repoName + getDataTugSectionForReadmeMD()
		_, _, err = c.client.Repositories.CreateFile(ctx, c.repoOwner, c.repoName, "README.md", &github.RepositoryContentFileOptions{
			Message: github.Ptr("feat: creates /README.md with ##DataTug section"),
			Content: []byte(newContent),
			Branch:  github.Ptr(c.branch),
		})
		if err != nil {
			return fmt.Errorf("failed to create root README.md: %w", err)
		}
	}

	c.stepRootReadme.Status = "updated."
	c.reportStatus()
	return nil
}

func (c *projectCreator) createProjectSummaryFile(pathToProjectFromRepoRoot string, project *datatug.Project) (err error) {
	projectFile := datatug.ProjectFile{
		ProjectItem: project.ProjectItem,
		Created: &datatug.ProjectCreated{
			At: time.Now().UTC(),
		},
	}

	filePath := path.Join(pathToProjectFromRepoRoot, filestore.ProjectSummaryFileName)

	return fmt.Errorf("createProjectFile is not implemented %v %v", filePath, projectFile)
}

func (c *projectCreator) createProjectReadme(pathToProjectFromRepoRoot string) error {
	filePath := path.Join(pathToProjectFromRepoRoot, "datatug/README.md")
	c.addEntry(filePath, "# DataTug Project\n\nThis directory contains DataTug project configuration.")
	return nil
}

//func (c *projectCreator) commitChanges(ctx context.Context, ref *github.Reference, tree *github.Tree) (err error) {
//	// 3. Create a commit
//	var parent *github.Commit
//	parent, _, err = c.client.Git.GetCommit(ctx, c.repoOwner, c.repoName, *ref.Object.SHA)
//	if err != nil {
//		err = fmt.Errorf("failed to get parent commit: %w", err)
//		return
//	}
//
//	var commit *github.Commit
//	commit, _, err = c.client.Git.CreateCommit(ctx, c.repoOwner, c.repoName, github.Commit{
//		Message: github.Ptr("chore: adds datatug project"),
//		Tree:    tree,
//		Parents: []*github.Commit{parent},
//	}, &github.CreateCommitOptions{})
//	if err != nil {
//		err = fmt.Errorf("failed to create commit: %w", err)
//		return
//	}
//
//	// 4. Update the reference
//	ref.Object.SHA = commit.SHA
//	_, _, err = c.client.Git.UpdateRef(ctx, c.repoOwner, c.repoName, ref.GetRef(), github.UpdateRef{
//		SHA:   commit.GetSHA(),
//		Force: github.Ptr(false),
//	})
//	return
//}

func (c *projectCreator) cloneRepo() (err error) {
	localDir := filestore.ExpandHome("projectPath")
	dirExists, _ := filestore.DirExists(localDir)
	if !dirExists {
		parent := filepath.Dir(localDir)
		_ = os.MkdirAll(parent, 0o755)
		cloneUrl := c.repo.GetCloneURL()
		if cloneUrl == "" {
			cloneUrl = fmt.Sprintf("https://github.com/%s/%s.git", c.repoOwner, c.repoName)
			_ = cloneUrl
		}
		//_, err = git.PlainClone(localDir, false, &git.CloneOptions{
		//	URL:      cloneUrl,
		//	Progress: NewTviewProgressWriter(tui, progressView),
		//})
		//if err != nil {
		//	tui.App.QueueUpdateDraw(func() {
		//		sneatnav.ShowErrorModal(tui, fmt.Errorf("failed to clone repository: %w", err))
		//	})
		//	return
		//}
	}
	return nil
}

func (c *projectCreator) addProjectToDataTugConfig(pathToProjectFromRepoRoot, projectTitle string) (err error) {
	// 4. Add project to DataTug app config
	projectID := fmt.Sprintf("github.com/%s/%s/%s", c.repoOwner, c.repoName, pathToProjectFromRepoRoot)
	projectDir := path.Join("~/datatug/", projectID)
	projectRef := dtconfig.ProjectRef{
		ID:    projectID,
		Path:  projectDir,
		Title: projectTitle,
	}
	if err = dtconfig.AddProjectToSettings(projectRef); err != nil {
		err = fmt.Errorf("failed to add repoName to DataTug app config: %w", err)
		return
	}
	return
}

// Helper to check for cancellation
func isCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
