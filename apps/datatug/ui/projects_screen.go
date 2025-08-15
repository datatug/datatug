package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
)

func newProjectsScreen(tui *tapp.TUI) tapp.Screen {
	return newDefaultLayout(tui, projectsRootScreen, func(tui *tapp.TUI) (tapp.Cell, error) {
		panel, err := newProjectsPanel(tui)
		return panel, err
	})
}

var _ tapp.Screen = (*projectsScreen)(nil)

type projectsScreen struct {
	tapp.ScreenBase
	//row *tapp.Row
}
