package gcloudui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GoCredentials(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
	menu := newMainMenu(tui, ScreenCredentials)

	list := tview.NewList()
	sneatv.SetPanelTitle(list.Box, "Google Cloud Projects")

	list.AddItem("Login", "", 'i', func() {})
	list.AddItem("Logout", "", 'o', func() {})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft, tcell.KeyEscape:
			tui.SetFocus(menu)
			return nil
		default:
			return event
		}
	})

	content := sneatnav.NewPanelFromList(tui, list)

	tui.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))
	return nil
}
