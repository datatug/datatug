package gcloudui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GoCredentials(cContext *GCloudContext, focusTo sneatnav.FocusTo) error {
	menu := newMainMenu(cContext, ScreenCredentials)

	list := tview.NewList()
	sneatv.SetPanelTitle(list.Box, "Google Cloud Projects")

	list.AddItem("Login", "", 'i', func() {})
	list.AddItem("Logout", "", 'o', func() {})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft, tcell.KeyEscape:
			cContext.TUI.SetFocus(menu)
			return nil
		default:
			return event
		}
	})

	content := sneatnav.NewPanelFromList(cContext.TUI, list)

	cContext.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))
	return nil
}
