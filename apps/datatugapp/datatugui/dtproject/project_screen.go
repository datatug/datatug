package dtproject

import (
	"fmt"
	"strings"

	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
)

func NewProjectPanel(tui *sneatnav.TUI, projectConfig *appconfig.ProjectConfig) sneatnav.Panel {
	content := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	projectTitle := GetProjectTitle(projectConfig)
	sneatv.SetPanelTitle(content.Box, fmt.Sprintf("Project: %s", projectTitle))
	return sneatnav.NewPanel(tui, sneatnav.WithBox(content, content.Box))
}

func GoProjectScreen(projectCtx ProjectContext) {
	tui := projectCtx.TUI()
	pConfig := projectCtx.Config()
	breadcrumbs := projectsBreadcrumbs(tui)
	title := GetProjectTitle(pConfig)
	if parts := strings.Split(title, "/"); len(parts) > 1 {
		title = parts[len(parts)-1]
	}
	breadcrumbs.Push(sneatv.NewBreadcrumb(title, nil))
	menu := newProjectMenuPanel(projectCtx, "project")
	content := NewProjectPanel(tui, pConfig)
	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToMenu))

	go func() {
		err := <-projectCtx.WatchProject()
		if err != nil {
			panic(fmt.Errorf("watch project error: %w", err))
		}
		project := projectCtx.Project()
		menu.SetProject(project)
	}()

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

func GetProjectShortTitle(p *appconfig.ProjectConfig) (projectTitle string) {
	projectTitle = GetProjectTitle(p)
	if parts := strings.Split(projectTitle, "/"); len(parts) > 1 {
		projectTitle = parts[len(parts)-1]
	}
	return projectTitle
}
