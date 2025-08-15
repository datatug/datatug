package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/rivo/tview"
)

func newDefaultLayout(tui *tapp.TUI, selectedMenuItem rootScreen, getContent func(tui *tapp.TUI) (tapp.Panel, error)) tapp.Screen {
	header := newHeaderPanel()

	grid := layoutGrid(header)

	addMainRow(tui, selectedMenuItem, grid, getContent)

	screen := &projectsScreen{
		ScreenBase: tapp.NewScreenBase(tui, grid, tapp.FullScreen()),
	}

	tui.SetRootScreen(screen)

	screen.TakeFocus()

	return screen
}

func addMainRow(tui *tapp.TUI, selectedMenuItem rootScreen, grid *tview.Grid, getContent func(tui *tapp.TUI) (tapp.Panel, error)) {
	menu := newHomeMenu(tui, selectedMenuItem)
	sidebar := newProjectsMenu(tui)

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
		AddItem(content, 1, 0, 1, 3, 0, 0, false).
		AddItem(sidebar, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.
		AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(content, 1, 1, 1, 1, 0, 100, false).
		AddItem(sidebar, 1, 2, 1, 1, 0, 100, false)

	grid.SetFocusFunc(func() {
		menu.TakeFocus()
	})

	_ = tapp.NewRow(tui.App,
		menu,
		content,
		sidebar,
	)
}
