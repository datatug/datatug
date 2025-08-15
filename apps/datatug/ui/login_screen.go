package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
)

func newLoginScreen(tui *tapp.TUI) tapp.Screen {
	return newDefaultLayout(tui, loginRootScreen, func(tui *tapp.TUI) (tapp.Cell, error) {
		panel, err := newLoginPanel(tui)
		return panel, err
	})
}

var _ tapp.Screen = (*loginScreen)(nil)

type loginScreen struct {
	tapp.ScreenBase
}
