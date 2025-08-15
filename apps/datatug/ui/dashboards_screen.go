package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-core/pkg/appconfig"
)

type dashboardsScreen struct {
	tapp.ScreenBase
}

func newDashboardsScreen(tui *tapp.TUI, project appconfig.ProjectConfig) tapp.Screen {
	main := newDashboardsPanel(project)

	_ = newDashboardsSidebar(tui)

	return &dashboardsScreen{
		ScreenBase: newProjectRootScreenBase(tui, project, ProjectScreenDashboards, main),
	}
}
