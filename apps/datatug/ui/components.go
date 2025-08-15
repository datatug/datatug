package ui

import (
	"github.com/rivo/tview"
)

func layoutGrid(header tview.Primitive) *tview.Grid {

	//footer := NewFooterPanel()

	grid := tview.NewGrid()

	grid. // Default grid settings
		SetRows(1, 0, 1).
		SetColumns(20, 0, 20).
		SetBorders(false)

	// Adds header and footer to the grid.
	grid.AddItem(header, 0, 0, 1, 3, 0, 0, false)
	//grid.AddItem(footer, 2, 0, 1, 3, 0, 0, false)

	return grid
}

func menuList() *tview.List {
	return tview.NewList().SetWrapAround(false)
}
