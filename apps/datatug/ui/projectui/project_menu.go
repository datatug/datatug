package projectui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/rivo/tview"
)

type ProjectScreenID string

const (
	ProjectScreenDashboards   = "dashboards"
	ProjectScreenEnvironments = "environments"
)

func newProjectMenuPanel(tui *sneatnav.TUI, project *appconfig.ProjectConfig, currentScreen ProjectScreenID) sneatnav.Panel {
	list := tview.NewList().
		//AddItem("Databases", "", 'D', nil).
		AddItem("Dashboards", "", 'B', func() {
			goProjectDashboards(tui, project)
		}).
		AddItem("Environments", "", 'E', func() {
			goEnvironmentsScreen(tui, project)
		})

	//AddItem("Queries", "", 'Q', nil).
	//AddItem("Web UI", "", 'W', nil)

	currentItem := -1
	switch currentScreen {
	case ProjectScreenDashboards:
		currentItem = 0
	case ProjectScreenEnvironments:
		currentItem = 1
	}
	if currentItem >= 0 {
		list.SetCurrentItem(currentItem)
	}

	sneatv.DefaultListStyle(list)
	sneatv.DefaultBorder(list.Box)

	return sneatnav.NewPanelFromList(tui, list)
}
