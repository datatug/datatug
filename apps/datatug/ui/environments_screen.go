package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug/packages/appconfig"
)

type environmentsScreen struct {
	tapp.ScreenBase
}

func newEnvironmentsScreen(tui *tapp.TUI, project appconfig.ProjectConfig) tapp.Screen {

	main := newEnvironmentsPanel(project)

	sidebar := newProjectsMenu(tui)

	return &environmentsScreen{
		ScreenBase: newProjectRootScreenBase(tui, project, ProjectScreenEnvironments, main, sidebar),
	}
}
