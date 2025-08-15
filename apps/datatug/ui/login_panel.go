package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/rivo/tview"
)

var _ tview.Primitive = (*loginPanel)(nil)
var _ tapp.Cell = (*loginPanel)(nil)

type loginPanel struct {
	tapp.PanelBase
}

func newLoginPanel(tui *tapp.TUI) (*loginPanel, error) {
	text := tview.NewTextView()
	text.SetText("Sign in browser: https://datatug.app")
	panel := &loginPanel{
		PanelBase: tapp.NewPanelBaseFromTextView(tui, text),
	}
	setPanelTitle(panel.PanelBase, "Login to DataTug")
	return panel, nil
}
