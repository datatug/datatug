package ui

import (
	"context"
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-cli/pkg/tvprimitives/breadcrumbs"
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/strongo/logus"
	"sort"
	"strconv"
)

var _ tview.Primitive = (*projectsPanel)(nil)
var _ tapp.Cell = (*projectsPanel)(nil)

type projectsPanel struct {
	tapp.PanelBase
	projects        []*appconfig.ProjectConfig
	selectProjectID string
	list            *tview.List
}

func goProjectsScreen(tui *tapp.TUI) error {
	content, err := getProjectsContent(tui)
	if err != nil {
		return err
	}
	tui.SetPanels(newDataTugMainMenu(tui, projectsRootScreen), content)
	tui.Header.Breadcrumbs.Clear()
	tui.Header.Breadcrumbs.Push(breadcrumbs.NewBreadcrumb("Projects", nil))
	return nil
}

func getProjectsContent(tui *tapp.TUI) (tapp.Panel, error) {
	panel, err := newProjectsPanel(tui)
	return panel, err
}

func newProjectsPanel(tui *tapp.TUI) (*projectsPanel, error) {
	list := tview.NewList()
	panel := &projectsPanel{
		PanelBase: tapp.NewPanelBaseFromList(tui, list),
		list:      list,
	}

	settings, err := appconfig.GetSettings()
	if err != nil {
		logus.Errorf(context.Background(), "Failed to get app settings: %v", err)
		//return nil, err
	}

	openProject := func(projectConfig appconfig.ProjectConfig) {
		projectScreen := newProjectScreen(tui, projectConfig)
		tui.PushScreen(projectScreen)
	}

	panel.projects = settings.Projects

	sort.Slice(panel.projects, func(i, j int) bool {
		return panel.projects[i].ID < panel.projects[j].ID
	})

	projectSelected := func(p *appconfig.ProjectConfig) {
		panel.selectProjectID = p.ID
		openProject(*p)
	}
	for i, p := range panel.projects {
		project := p
		list.AddItem(project.ID, project.Url, rune(strconv.Itoa(i + 1)[0]), func() {
			projectSelected(project)
		})
	}

	defaultListStyle(list)

	setPanelTitle(panel.PanelBase, "Projects")

	return panel, nil
}

func (p *projectsPanel) Draw(screen tcell.Screen) {
	var selectedItem = -1

	for i, proj := range p.projects {
		if proj.ID == p.selectProjectID {
			selectedItem = i
		}
	}
	if selectedItem >= 0 {
		p.list.SetCurrentItem(selectedItem)
	}
	p.list.Draw(screen)
}
