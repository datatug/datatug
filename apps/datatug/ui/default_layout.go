package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-cli/pkg/tvprimitives/breadcrumbs"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newDefaultLayout(tui *tapp.TUI, selectedMenuItem rootScreen, getContent func(tui *tapp.TUI) (tapp.Panel, error)) (
	tapp.Screen, *breadcrumbs.Breadcrumbs,
) {

	grid := layoutGrid(tui.Header)

	addMainRow(tui, selectedMenuItem, grid, getContent, tui.Header.Breadcrumbs)

	screen := &projectsScreen{
		ScreenBase: tapp.NewScreenBase(tui, grid, tapp.FullScreen()),
	}

	tui.SetRootScreen(screen)

	screen.TakeFocus()

	return screen, tui.Header.Breadcrumbs
}

func addMainRow(tui *tapp.TUI, selectedMenuItem rootScreen, grid *tview.Grid, getContent func(tui *tapp.TUI) (tapp.Panel, error), header *breadcrumbs.Breadcrumbs) {
	menu := newDataTugMainMenu(tui, selectedMenuItem)

	content, err := getContent(tui)
	if err != nil {
		panic(err)
	}
	if content == nil {
		panic("getContent() returned nil")
	}

	// Layout for screens narrower than 100 cells (menu and sidebar are hidden).
	grid.
		AddItem(menu, 0, 0, 0, 0, 0, 0, false).
		AddItem(content, 1, 0, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.
		AddItem(menu, 1, 0, 1, 1, 0, 100, true).
		AddItem(content, 1, 1, 1, 1, 0, 100, false)

	// When header wants to move to the next row's first cell, that is our menu.
	header.SetNextFocusTarget(menu)

	// Allow keyboard navigation from the menu to the header with Shift+Tab (Backtab) or Up arrow.
	// This enables Breadcrumbs to receive focus and thus its InputHandler to be called.
	menu.list.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		// Handle the logic from newDataTugMainMenu: move focus to breadcrumbs when on first item
		if ev.Key() == tcell.KeyUp {
			if menu.list.GetCurrentItem() == 0 {
				tui.App.SetFocus(tui.Header.Breadcrumbs)
				return nil
			}
		}
		if ev.Key() == tcell.KeyBacktab {
			// Move focus to header (breadcrumbs) when Shift+Tab or Up arrow is pressed on the menu.
			tui.App.SetFocus(header)
			return nil // consume the event
		}
		return ev
	})

	grid.SetFocusFunc(func() {
		menu.TakeFocus()
	})

	_ = tapp.NewRow(tui.App,
		menu,
		content,
	)
}
