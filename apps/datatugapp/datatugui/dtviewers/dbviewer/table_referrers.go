package dbviewer

import (
	"fmt"

	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewReferrersBox(_ dtviewers.CollectionContext) *tview.Table {
	table := tview.NewTable()
	table.SetTitle("Referrers")
	table.SetBorder(true)
	table.SetBorderColor(tcell.ColorDarkSlateGray)
	//table.SetSelectable(true, false)

	newRow := func(row int, referrer, fkCols string) {
		table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%s(%s)", referrer, fkCols)))
		return
	}

	newRow(0, "Orders", "CustomerID")
	return table
}
