package ui

import (
	tapp2 "github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug/packages/appconfig"
)

func newProjectScreen(tui *tapp2.TUI, project appconfig.ProjectConfig) tapp2.Screen {
	return newEnvironmentsScreen(tui, project)
}
