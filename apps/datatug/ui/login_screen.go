package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
)

func newLoginScreen(tui *tapp.TUI) tapp.Screen {
	screen, _ := newDefaultLayout(tui, loginRootScreen, func(tui *tapp.TUI) (tapp.Panel, error) {
		panel, err := newLoginPanel(tui)
		return panel, err
	})
	return screen
}

var _ tapp.Screen = (*loginScreen)(nil)

type loginScreen struct {
	tapp.ScreenBase
}
