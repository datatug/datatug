package sneatnav

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ShowErrorModal(tui *TUI, err error) {
	text := tview.NewTextView()
	text.SetText(err.Error()).SetTextColor(tcell.ColorRed)
	content := NewPanel(tui, WithBox(text, text.Box))
	tui.SetPanels(tui.Menu, content)
}
