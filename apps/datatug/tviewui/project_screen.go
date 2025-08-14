package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug/packages/appconfig"
)

func newProjectScreen(tui *tapp.TUI, project appconfig.ProjectConfig) tapp.Screen {
	return newEnvironmentsScreen(tui, project)
}
