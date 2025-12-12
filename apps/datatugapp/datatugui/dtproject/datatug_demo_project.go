package dtproject

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	git "github.com/go-git/go-git/v5"
	"github.com/rivo/tview"
)

const datatugDemoProjectID = "datatug-demo-project"
const datatugDemoProjectGitHubRepoID = "datatug/datatug-demo-project"
const datatugDemoProjectGitHubRepoURL = "https://github.com/" + datatugDemoProjectGitHubRepoID
const datatugDemoProjectDir = "~/datatug/" + datatugDemoProjectID

func openDatatugDemoProject(tui *sneatnav.TUI) {
	// Expand home in path like ~/...
	projectDir := expandHome(datatugDemoProjectDir)

	projectDirExists, err := dirExists(projectDir)
	if err != nil {
		panic(err)
	}
	if !projectDirExists {

		progressText := tview.NewTextView()
		progressText.SetTitle("Cloning project...")
		progressPanel := sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(progressText, progressText.Box))
		tui.SetPanels(tui.Menu, progressPanel, sneatnav.WithFocusTo(sneatnav.FocusToContent))

		go func() {
			// Ensure parent directory exists
			parent := filepath.Dir(projectDir)
			if err = os.MkdirAll(parent, 0o755); err != nil {
				panic(err)
			}
			// Clone public GitHub repository datatugDemoProjectGitHubRepoID into datatugDemoProjectDir using go-git
			if _, err = git.PlainClone(projectDir, false, &git.CloneOptions{
				URL:      datatugDemoProjectGitHubRepoURL,
				Progress: newTviewProgressWriter(tui, progressText),
				// Depth: 1, // uncomment for shallow clone if desired
			}); err != nil {
				panic("git clone failed: " + err.Error())
			}
			tui.App.QueueUpdateDraw(func() {
				p := &appconfig.ProjectConfig{
					ID:  datatugDemoProjectID,
					Url: datatugDemoProjectGitHubRepoURL,
				}
				GoProjectScreen(tui, p)
			})
		}()
	}
}
func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err // some other error
	}
	return info.IsDir(), nil
}

// expandHome expands leading ~ to the user's home directory.
func expandHome(p string) string {
	if p == "" {
		return p
	}
	if strings.HasPrefix(p, "~/") || p == "~" {
		home, err := os.UserHomeDir()
		if err == nil {
			if p == "~" {
				return home
			}
			return filepath.Join(home, strings.TrimPrefix(p, "~/"))
		}
	}
	return p
}

// tviewProgressWriter implements io.Writer and appends text to a TextView safely via tview.Application.
type tviewProgressWriter struct {
	tui *sneatnav.TUI
	tv  *tview.TextView
}

func newTviewProgressWriter(tui *sneatnav.TUI, tv *tview.TextView) *tviewProgressWriter {
	return &tviewProgressWriter{tui: tui, tv: tv}
}

func (w *tviewProgressWriter) Write(p []byte) (n int, err error) {
	// Ensure UI updates happen on the application goroutine
	w.tui.App.QueueUpdateDraw(func() {
		w.tv.SetText(string(p))
	})
	return len(p), nil
}
