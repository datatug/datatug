package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-cli/pkg/tvprimitives/breadcrumbs"
)

func newLoginScreen(tui *tapp.TUI) tapp.Screen {
	screen, _ := newDefaultLayout(tui, loginRootScreen, func(tui *tapp.TUI) (tapp.Panel, error) {
		panel, err := newLoginPanel(tui)
		return panel, err
	})
	tui.Header.Breadcrumbs.Clear()
	tui.Header.Breadcrumbs.Push(breadcrumbs.NewBreadcrumb("Login", nil))
	return screen
}

var _ tapp.Screen = (*loginScreen)(nil)

type loginScreen struct {
	tapp.ScreenBase
}
