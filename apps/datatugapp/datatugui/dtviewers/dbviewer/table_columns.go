package dbviewer

import (
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewColumnsBox(_ dtviewers.CollectionContext) *tview.Table {
	table := tview.NewTable()
	table.SetTitle(`[gray]Categories:[-] Columns`)
	sneatv.DefaultBorderWithoutPadding(table.Box)

	newRow := func(row int, name, dbType string) {
		table.SetCell(row, 0, tview.NewTableCell(name))
		table.SetCell(row, 1, tview.NewTableCell(dbType).SetTextColor(tcell.ColorGray))
		return
	}
	newRow(0, "CategoryID", "INT")
	newRow(1, "CategoryName", "TEXT")
	newRow(2, "Description", "TEXT")
	newRow(3, "Picture", "BLOB")
	return table
}
