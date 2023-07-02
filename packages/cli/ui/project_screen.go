package ui

import (
	"github.com/datatug/datatug/packages/cli/config"
	"github.com/datatug/datatug/packages/cli/tapp"
	"github.com/rivo/tview"
)

type projectScreen struct {
	tapp.ScreenBase
}

func NewProjectScreen(tui *tapp.TUI, project config.ProjectConfig) tapp.Screen {

	menu := newProjectMenu(tui)

	sideBar := newProjectsMenu(tui)

	header := newHeaderPanel(tui, project.ID)

	footer := NewFooterPanel()

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(20, 0, 20).
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(footer, 2, 0, 1, 3, 0, 0, false).
		AddItem(sideBar, 1, 0, 1, 1, 0, 0, false)

	main := newProjectPanel(tui, project)

	// Layout for screens narrower than 100 cells (menu and sidebar are hidden).
	grid.
		AddItem(menu, 0, 0, 0, 0, 0, 0, false).
		AddItem(main, 1, 0, 1, 3, 0, 0, false).
		AddItem(sideBar, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(main, 1, 1, 1, 1, 0, 100, false).
		AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	grid.SetFocusFunc(func() {
		menu.TakeFocus()
	})

	_ = tapp.NewRow(tui.App, menu, main, sideBar)

	screen := &projectScreen{
		ScreenBase: tapp.NewScreenBase(tui, grid, tapp.FullScreen()),
	}

	screen.TakeFocus()

	return screen
}
