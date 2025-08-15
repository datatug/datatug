package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/rivo/tview"
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

func newLoginPanel(tui *tapp.TUI) (*loginPanel, error) {
	textView := tview.NewTextView().SetText(string("Login to DataTug"))
	panel := &loginPanel{
		PanelBase: tapp.NewPanelBase(tui, textView, textView.Box),
	}
	return panel, nil
}
