package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-core/pkg/appconfig"
)

type dashboardsScreen struct {
	sneatnav.ScreenBase
}

func newDashboardsScreen(tui *sneatnav.TUI, project appconfig.ProjectConfig) sneatnav.Screen {
	main := newDashboardsPanel(project)

	_ = newDashboardsSidebar(tui)

	return &dashboardsScreen{
		ScreenBase: newProjectRootScreenBase(tui, project, ProjectScreenDashboards, main),
	}
}
