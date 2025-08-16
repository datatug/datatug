package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/rivo/tview"
)

func newDefaultLayout(
	tui *sneatnav.TUI, selectedMenuItem rootScreen, getContent func(tui *sneatnav.TUI) (sneatnav.Panel, error),
) sneatnav.Screen {
	addMainRow(tui, selectedMenuItem, tui.Grid, getContent)

	return nil
}

func addMainRow(
	tui *sneatnav.TUI, selectedMenuItem rootScreen, grid *tview.Grid,
	getContent func(tui *sneatnav.TUI) (sneatnav.Panel, error),
) {
	menu := newDataTugMainMenu(tui, selectedMenuItem)

	content, err := getContent(tui)
	if err != nil {
		panic(err)
	}
	if content == nil {
		panic("getContent() returned nil")
	}

	// Allow keyboard navigation from the menu to the header with Shift+Tab (Backtab) or Up arrow.
	// This enables Breadcrumbs to receive focus and thus its InputHandler to be called.

	grid.SetFocusFunc(func() {
		menu.TakeFocus()
	})

	_ = sneatnav.NewRow(tui.App,
		menu,
		content,
	)
}
