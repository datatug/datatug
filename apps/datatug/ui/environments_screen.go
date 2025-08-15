package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-core/pkg/appconfig"
)

type environmentsScreen struct {
	tapp.ScreenBase
}

func newEnvironmentsScreen(tui *tapp.TUI, project appconfig.ProjectConfig) tapp.Screen {

	main := newEnvironmentsPanel(project)

	_ = newProjectsMenu(tui)

	return &environmentsScreen{
		ScreenBase: newProjectRootScreenBase(tui, project, ProjectScreenEnvironments, main),
	}
}
