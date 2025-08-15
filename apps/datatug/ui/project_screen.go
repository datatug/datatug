package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-core/pkg/appconfig"
)

func newProjectScreen(tui *tapp.TUI, project appconfig.ProjectConfig) tapp.Screen {
	return newEnvironmentsScreen(tui, project)
}
