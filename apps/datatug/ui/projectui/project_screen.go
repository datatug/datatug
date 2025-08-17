package projectui

import (
	"fmt"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/rivo/tview"
	"strings"
)

func NewProjectPanel(tui *sneatnav.TUI, projectConfig *appconfig.ProjectConfig) sneatnav.Panel {
	content := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	projectTitle := GetProjectTitle(projectConfig)
	sneatv.SetPanelTitle(content.Box, fmt.Sprintf("Project: %s", projectTitle))

	sneatv.DefaultBorder(content.Box)

	return sneatnav.NewPanelFromTextView(tui, content)
}

func GoProjectScreen(tui *sneatnav.TUI, p *appconfig.ProjectConfig) {
	menu := newProjectMenuPanel(tui, p, "project")
	content := NewProjectPanel(tui, p)
	tui.Header.Breadcrumbs().Clear()
	tui.Header.Breadcrumbs().Push(sneatv.NewBreadcrumb("Projects", nil))
	title := GetProjectTitle(p)
	if parts := strings.Split(title, "/"); len(parts) > 1 {
		title = parts[len(parts)-1]
	}
	tui.Header.Breadcrumbs().Push(sneatv.NewBreadcrumb(title, nil))
	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToMenu))
}

func GetProjectTitle(p *appconfig.ProjectConfig) (projectTitle string) {
	projectTitle = p.Title
	if projectTitle == "" {
		projectTitle = p.ID
	}
	if projectTitle == "" {
		projectTitle = p.Url
	}
	return projectTitle
}
