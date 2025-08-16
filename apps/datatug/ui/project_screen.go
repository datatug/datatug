package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-core/pkg/appconfig"
)

func newProjectScreen(tui *sneatnav.TUI, project appconfig.ProjectConfig) sneatnav.Screen {
	return newEnvironmentsScreen(tui, project)
}
