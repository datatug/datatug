package sneatnav

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func MainMenuList(tui *TUI) *tview.List {
	list := tview.NewList()
	list.SetWrapAround(false)
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			tui.SetFocus(tui.Menu)
			return nil
		case tcell.KeyRight:
			tui.SetFocus(tui.Content)
			return nil
		default:
			return event
		}
	})
	return list
}
