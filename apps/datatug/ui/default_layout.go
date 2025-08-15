package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-cli/pkg/tvprimitives/breadcrumbs"
	"github.com/rivo/tview"
)

func newDefaultLayout(
	tui *tapp.TUI, selectedMenuItem rootScreen, getContent func(tui *tapp.TUI) (tapp.Panel, error),
) (
	tapp.Screen, *breadcrumbs.Breadcrumbs,
) {

	addMainRow(tui, selectedMenuItem, tui.Grid, getContent, tui.Header.Breadcrumbs)

	return nil, tui.Header.Breadcrumbs
}

func addMainRow(
	tui *tapp.TUI, selectedMenuItem rootScreen, grid *tview.Grid,
	getContent func(tui *tapp.TUI) (tapp.Panel, error),
	header *breadcrumbs.Breadcrumbs,
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

	_ = tapp.NewRow(tui.App,
		menu,
		content,
	)
}
