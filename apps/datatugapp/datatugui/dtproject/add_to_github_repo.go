package dtproject

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/datatug/datatug-core/pkg/storage/filestore"
	"github.com/datatug/datatug/pkg/auth/ghauth"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/gdamore/tcell/v2"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v80/github"
	"github.com/rivo/tview"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

/*
## High level flow:
Go CLI
 ├─ Request device code from GitHub
 ├─ Display user code and verification URL
 ├─ Poll GitHub for access_token
 ├─ CLI lists user repositories
 ├─ User selects repo
 ├─ CLI creates datatug/README.md via GitHub API

Uses https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#device-flow
*/

// ShowAddToGitHubRepo
// - gets user credentials for GitHub api via OAuth2 Device Flow
// - selects repository to be used to store datatug project.
// - adds a `datatug` directory with config and README.md to the root of an existing GitHub repo.
// - adds a 'DataTug' section to the root README.md files linked to the `datatug` directory.
func ShowAddToGitHubRepo(tui *sneatnav.TUI) {
	ctx := context.Background()

	clientID := "Ov23liAIKfguW2oYiore"
	clientSecret := os.Getenv("GITHUB_OAUTH_SECRET")

	var startAuth func()

	reauth := func() {
		_ = ghauth.DeleteToken()
		startAuth()
	}

	useToken := func(token *oauth2.Token) {
		ts := oauth2.StaticTokenSource(token)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		// List repositories
		repos, _, err := client.Repositories.ListByAuthenticatedUser(ctx, nil)
		if err != nil {
			// If listing repositories fails, the token might be invalid or expired.
			// In that case, we might want to delete it and restart auth.
			_ = ghauth.DeleteToken()
			startAuth()
			return
		}

		showRepoSelection(tui, client, repos, reauth)
	}

	startAuth = func() {
		go func() {
			deviceRes, err := ghauth.RequestDeviceCode(ctx, clientID)
			if err != nil {
				tui.App.QueueUpdateDraw(func() {
					sneatnav.ShowErrorModal(tui, fmt.Errorf("failed to request device code: %w", err))
				})
				return
			}

			tui.App.QueueUpdateDraw(func() {
				statusText := tview.NewTextView().
					SetDynamicColors(true).
					SetTextAlign(tview.AlignCenter).
					SetText(fmt.Sprintf("\nGo to %s\n\nEnter code: [yellow]%s[-]\n\nWaiting for authorization...", deviceRes.VerificationURI, deviceRes.UserCode))

				copyMessage := tview.NewTextView().
					SetTextAlign(tview.AlignCenter).
					SetTextColor(tcell.ColorGreen)

				form := tview.NewForm().
					AddButton("Copy Code", func() {
						_ = clipboard.WriteAll(deviceRes.UserCode)
						copyMessage.SetText("The code has been copied to clipboard.")
						go func() {
							time.Sleep(2 * time.Second)
							tui.App.QueueUpdateDraw(func() {
								copyMessage.SetText("")
							})
						}()
					}).
					AddButton("Cancel", func() {
						_ = GoProjectsScreen(tui, sneatnav.FocusToContent)
					})
				form.SetButtonsAlign(tview.AlignCenter).
					SetButtonBackgroundColor(tcell.ColorDarkBlue).
					SetButtonTextColor(tcell.ColorWhite)

				flex := tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(statusText, 0, 1, false).
					AddItem(copyMessage, 1, 0, false).
					AddItem(form, 3, 1, true)

				flex.SetBorder(true).SetTitle("GitHub Device Activation")

				panel := sneatnav.NewPanel(tui, sneatnav.WithBox(flex, flex.Box))
				tui.SetPanels(nil, panel)

				// Update polling message
				updateStatus := func(attempt int) {
					tui.App.QueueUpdateDraw(func() {
						statusText.SetText(fmt.Sprintf("\nGo to %s\n\nEnter code: [yellow]%s[-]\n\nWaiting for authorization (attempt %d)...", deviceRes.VerificationURI, deviceRes.UserCode, attempt))
					})
				}

				go func() {
					token, err := ghauth.PollForToken(ctx, clientID, clientSecret, deviceRes.DeviceCode, deviceRes.Interval, updateStatus)
					tui.App.QueueUpdateDraw(func() {
						if err != nil {
							sneatnav.ShowErrorModal(tui, fmt.Errorf("authentication failed: %w", err))
							return
						}

						if err := ghauth.SaveToken(token); err != nil {
							// Log error but proceed
							fmt.Printf("failed to save token: %v\n", err)
						}

						useToken(token)
					})
				}()
			})
		}()
	}

	// Try to get token from keyring
	if token, err := ghauth.GetToken(); err == nil && token != nil {
		useToken(token)
	} else {
		startAuth()
	}
}

func showRepoSelection(tui *sneatnav.TUI, client *github.Client, repos []*github.Repository, reauth func()) {
	tree := tview.NewTreeView()
	tree.SetTitle("Select GitHub Repository").SetBorder(true)

	root := tview.NewTreeNode("Repositories").SetSelectable(false)
	tree.SetRoot(root)

	// Group repos by owner
	reposByOwner := make(map[string][]*github.Repository)
	var owners []string
	for _, repo := range repos {
		owner := repo.GetOwner().GetLogin()
		if _, ok := reposByOwner[owner]; !ok {
			owners = append(owners, owner)
		}
		reposByOwner[owner] = append(reposByOwner[owner], repo)
	}
	sort.Strings(owners)

	for _, owner := range owners {
		ownerNode := tview.NewTreeNode(owner).
			SetColor(tcell.ColorLightBlue).
			SetSelectable(true).
			SetExpanded(false)
		root.AddChild(ownerNode)

		ownerRepos := reposByOwner[owner]
		sort.Slice(ownerRepos, func(i, j int) bool {
			return ownerRepos[i].GetName() < ownerRepos[j].GetName()
		})

		for _, repo := range ownerRepos {
			r := repo
			repoNode := tview.NewTreeNode(repo.GetName()).
				SetReference(r).
				SetSelectedFunc(func() {
					AddToGitHubRepo(tui, client, r, repos, reauth)
				})
			ownerNode.AddChild(repoNode)
		}
	}

	cancelNode := tview.NewTreeNode("Cancel").
		SetReference("cancel").
		SetColor(tcell.ColorRed).
		SetSelectedFunc(func() {
			_ = GoProjectsScreen(tui, sneatnav.FocusToContent)
		})
	root.AddChild(cancelNode)

	reauthNode := tview.NewTreeNode("Re-authenticate").
		SetReference("reauth").
		SetColor(tcell.ColorYellow).
		SetSelectedFunc(func() {
			reauth()
		})
	root.AddChild(reauthNode)

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentNode := tree.GetCurrentNode()
		if currentNode == nil {
			return event
		}

		switch event.Key() {
		case tcell.KeyEnter, tcell.KeyRune:
			if event.Key() == tcell.KeyRune && event.Rune() != ' ' {
				return event
			}
			// If it's an owner node (has children and no Repository reference), toggle expansion
			if len(currentNode.GetChildren()) > 0 && currentNode.GetReference() == nil {
				currentNode.SetExpanded(!currentNode.IsExpanded())
				return nil
			}
			return event
		case tcell.KeyLeft:
			if currentNode.IsExpanded() {
				currentNode.SetExpanded(false)
				return nil
			}
			return event
		default:
			return event
		}
	})

	if len(root.GetChildren()) > 0 {
		tree.SetCurrentNode(root.GetChildren()[0])
	}

	panel := sneatnav.NewPanel(tui, sneatnav.WithBox(tree, tree.Box))
	tui.SetPanels(nil, panel)
}

func AddToGitHubRepo(tui *sneatnav.TUI, client *github.Client, repo *github.Repository, repos []*github.Repository, reauth func()) {
	owner := repo.GetOwner().GetLogin()
	name := repo.GetName()
	branch := repo.GetDefaultBranch()

	projectID := fmt.Sprintf("github.com/%s/%s", owner, name)
	projectTitle := fmt.Sprintf("%s @ github.com/%s", name, owner)
	projectDir := "~/datatug/" + projectID

	// UI for progress
	progressView := tview.NewTextView().SetDynamicColors(true)
	progressView.SetBorder(true).SetTitle("Setting up DataTug in " + repo.GetFullName())

	ctx, cancel := context.WithCancel(context.Background())

	cancelButton := tview.NewButton("Cancel").SetSelectedFunc(func() {
		cancel()
		if repos != nil {
			showRepoSelection(tui, client, repos, reauth)
		} else {
			_ = GoProjectsScreen(tui, sneatnav.FocusToContent)
		}
	})

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(progressView, 0, 1, false).
		AddItem(cancelButton, 1, 0, true)

	panel := sneatnav.NewPanel(tui, sneatnav.WithBox(layout, layout.Box))
	tui.SetPanels(nil, panel)

	steps := []string{
		"Add DataTug project files to repository",
		"Add DataTug section to /README.md",
		"Cloning project repository to " + projectDir,
		"Add project to DataTug app config",
	}

	updateProgress := func(currentStep int, status string) {
		tui.App.QueueUpdateDraw(func() {
			var sb strings.Builder
			for i, step := range steps {
				if i < currentStep {
					sb.WriteString(fmt.Sprintf("- %s - [green]done[-]\n", step))
				} else if i == currentStep {
					sb.WriteString(fmt.Sprintf("- %s - [yellow]%s[-]\n", step, status))
				} else {
					sb.WriteString(fmt.Sprintf("- %s\n", step))
				}
			}
			progressView.SetText(sb.String())
		})
	}

	go func() {
		// Helper to check for cancellation
		isCancelled := func() bool {
			select {
			case <-ctx.Done():
				return true
			default:
				return false
			}
		}

		// 1. Create datatug directory with config and README.md in a single commit
		updateProgress(0, "creating...")
		if isCancelled() {
			return
		}

		configContent := `{
  "id": "` + name + `",
  "title": "` + name + `"
}`
		configFilePath := "datatug/" + filestore.ProjectSummaryFileName
		readmeContent := "# DataTug Project\n\nThis directory contains DataTug project configuration."
		readmeFilePath := "datatug/README.md"

		// We use the Git Data API to create multiple files in a single commit.
		// 1. Get the latest commit of the branch
		ref, _, err := client.Git.GetRef(ctx, owner, name, "heads/"+branch)
		if err != nil {
			// If repository is empty, we need to create the first commit
			if gerr, ok := err.(*github.ErrorResponse); ok && (gerr.Response.StatusCode == 404 || gerr.Response.StatusCode == 409) {
				// Create initial README.md to initialize the repository
				updateProgress(0, "initializing repository...")
				_, _, err = client.Repositories.CreateFile(ctx, owner, name, "README.md", &github.RepositoryContentFileOptions{
					Message: github.Ptr("feat: initial commit"),
					Content: []byte("# " + name + "\n\nDataTug project repository."),
					Branch:  github.Ptr(branch),
				})
				if err != nil {
					tui.App.QueueUpdateDraw(func() {
						sneatnav.ShowErrorModal(tui, fmt.Errorf("failed to initialize repository: %w", err))
					})
					return
				}
				// Retry getting the ref
				ref, _, err = client.Git.GetRef(ctx, owner, name, "heads/"+branch)
			}
			if err != nil {
				tui.App.QueueUpdateDraw(func() {
					sneatnav.ShowErrorModal(tui, fmt.Errorf("failed to get branch ref: %w", err))
				})
				return
			}
		}

		// 2. Create a tree with the new files
		entries := []*github.TreeEntry{
			{
				Path:    github.Ptr(configFilePath),
				Type:    github.Ptr("blob"),
				Mode:    github.Ptr("100644"),
				Content: github.Ptr(configContent),
			},
			{
				Path:    github.Ptr(readmeFilePath),
				Type:    github.Ptr("blob"),
				Mode:    github.Ptr("100644"),
				Content: github.Ptr(readmeContent),
			},
		}

		// Check if files already exist to avoid overwriting or redundant commits
		// Actually, if we just want to ensure they exist, we can check first.
		existingConfig, _, _, _ := client.Repositories.GetContents(ctx, owner, name, configFilePath, &github.RepositoryContentGetOptions{Ref: branch})
		existingReadme, _, _, _ := client.Repositories.GetContents(ctx, owner, name, readmeFilePath, &github.RepositoryContentGetOptions{Ref: branch})

		var entriesToCreate []*github.TreeEntry
		if existingConfig == nil {
			entriesToCreate = append(entriesToCreate, entries[0])
		}
		if existingReadme == nil {
			entriesToCreate = append(entriesToCreate, entries[1])
		}

		if len(entriesToCreate) > 0 {
			tree, _, err := client.Git.CreateTree(ctx, owner, name, *ref.Object.SHA, entriesToCreate)
			if err != nil {
				tui.App.QueueUpdateDraw(func() {
					sneatnav.ShowErrorModal(tui, fmt.Errorf("failed to create tree: %w", err))
				})
				return
			}

			// 3. Create a commit
			parent, _, err := client.Git.GetCommit(ctx, owner, name, *ref.Object.SHA)
			if err != nil {
				tui.App.QueueUpdateDraw(func() {
					sneatnav.ShowErrorModal(tui, fmt.Errorf("failed to get parent commit: %w", err))
				})
				return
			}

			commit, _, err := client.Git.CreateCommit(ctx, owner, name, github.Commit{
				Message: github.Ptr("chore: add datatug project"),
				Tree:    tree,
				Parents: []*github.Commit{parent},
			}, &github.CreateCommitOptions{})
			if err != nil {
				tui.App.QueueUpdateDraw(func() {
					sneatnav.ShowErrorModal(tui, fmt.Errorf("failed to create commit: %w", err))
				})
				return
			}

			// 4. Update the reference
			ref.Object.SHA = commit.SHA
			_, _, err = client.Git.UpdateRef(ctx, owner, name, ref.GetRef(), github.UpdateRef{
				SHA:   commit.GetSHA(),
				Force: github.Ptr(false),
			})
			if err != nil {
				tui.App.QueueUpdateDraw(func() {
					sneatnav.ShowErrorModal(tui, fmt.Errorf("failed to update ref: %w", err))
				})
				return
			}
		}

		if isCancelled() {
			return
		}

		// 2. Add 'DataTug' section to root README.md
		updateProgress(1, "updating...")
		rootReadme, _, err := client.Repositories.GetReadme(ctx, owner, name, &github.RepositoryContentGetOptions{Ref: branch})
		if err == nil {
			content, _ := rootReadme.GetContent()
			if !strings.Contains(content, "## DataTug") {
				newContent := content + "\n\n## DataTug\n\nThis project is enhanced with [DataTug](https://datatug.app). See the [datatug](./datatug) directory for details.\n"
				_, _, err = client.Repositories.UpdateFile(ctx, owner, name, rootReadme.GetPath(), &github.RepositoryContentFileOptions{
					Message: github.Ptr("chore: add DataTug section to README"),
					Content: []byte(newContent),
					SHA:     rootReadme.SHA,
					Branch:  github.Ptr(branch),
				})
				if err != nil {
					tui.App.QueueUpdateDraw(func() {
						sneatnav.ShowErrorModal(tui, fmt.Errorf("failed to update root README.md: %w", err))
					})
					return
				}
			}
		} else {
			newContent := "# " + name + "\n\n## DataTug\n\nThis project is enhanced with [DataTug](https://datatug.app). See the [datatug](./datatug) directory for details.\n"
			_, _, err = client.Repositories.CreateFile(ctx, owner, name, "README.md", &github.RepositoryContentFileOptions{
				Message: github.Ptr("feat: add README with DataTug section"),
				Content: []byte(newContent),
				Branch:  github.Ptr(branch),
			})
			if err != nil {
				tui.App.QueueUpdateDraw(func() {
					sneatnav.ShowErrorModal(tui, fmt.Errorf("failed to create root README.md: %w", err))
				})
				return
			}
		}

		if isCancelled() {
			return
		}

		// 3. Cloning project repository
		updateProgress(3, "cloning...")
		localDir := filestore.ExpandHome(projectDir)
		dirExists, _ := filestore.DirExists(localDir)
		if !dirExists {
			parent := filepath.Dir(localDir)
			_ = os.MkdirAll(parent, 0o755)
			cloneUrl := repo.GetCloneURL()
			if cloneUrl == "" {
				cloneUrl = fmt.Sprintf("https://github.com/%s/%s.git", owner, name)
			}
			_, err = git.PlainClone(localDir, false, &git.CloneOptions{
				URL:      cloneUrl,
				Progress: NewTviewProgressWriter(tui, progressView),
			})
			if err != nil {
				tui.App.QueueUpdateDraw(func() {
					sneatnav.ShowErrorModal(tui, fmt.Errorf("failed to clone repository: %w", err))
				})
				return
			}
		}

		if isCancelled() {
			return
		}

		// 4. Add project to DataTug app config
		updateProgress(4, "updating...")
		if err := AddProjectToSettings(projectID, projectTitle, projectDir); err != nil {
			tui.App.QueueUpdateDraw(func() {
				sneatnav.ShowErrorModal(tui, fmt.Errorf("failed to add repo to DataTug app config: %w", err))
			})
			return
		}

		updateProgress(5, "") // All done

		time.Sleep(500 * time.Millisecond)

		if isCancelled() {
			return
		}

		tui.App.QueueUpdateDraw(func() {
			loader := filestore.NewProjectsLoader("~/datatug")
			pConfig := &appconfig.ProjectConfig{
				ID:    projectID,
				Title: projectTitle,
				Path:  projectID,
			}
			projectCtx := NewProjectContext(tui, pConfig, loader)
			GoProjectScreen(projectCtx)
		})
	}()
}

func AddProjectToSettings(id, title, path string) error {
	settings, err := appconfig.GetSettings()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to get DataTug app settings: %w", err)
	}

	// Check if already exists
	var project *appconfig.ProjectConfig
	for _, p := range settings.Projects {
		if p.ID == id {
			project = p
			break
		}
	}

	if project == nil {
		project = &appconfig.ProjectConfig{ID: id}
		settings.Projects = append(settings.Projects, project)
	}
	project.Title = title
	project.Path = path

	return SaveDataTugAppSettings(settings)
}

func SaveDataTugAppSettings(settings appconfig.Settings) error {
	configFilePath := appconfig.GetConfigFilePath()
	f, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to create settings file: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	if settings.Server != nil && settings.Server.IsEmpty() {
		settings.Server = nil
	}
	if settings.Client != nil && settings.Client.IsEmpty() {
		settings.Client = nil
	}

	encoder := yaml.NewEncoder(f)
	if err := encoder.Encode(settings); err != nil {
		return fmt.Errorf("failed to encode settings: %w", err)
	}
	return nil
}
