package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-core/pkg/appconfig"
)

type environmentsScreen struct {
	sneatnav.ScreenBase
}

func newEnvironmentsScreen(tui *sneatnav.TUI, project appconfig.ProjectConfig) sneatnav.Screen {

	main := newEnvironmentsPanel(project)

	_ = newProjectsMenu(tui)

	return &environmentsScreen{
		ScreenBase: newProjectRootScreenBase(tui, project, ProjectScreenEnvironments, main),
	}
}
