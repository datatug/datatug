package sneatnav

import "github.com/rivo/tview"

func ShowErrorModal(tui *TUI, err error) {
	text := tview.NewTextView()
	text.SetText(err.Error())
	NewPanel(tui, WithBox(text, text.Box))
}
