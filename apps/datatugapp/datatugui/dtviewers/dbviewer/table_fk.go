package dbviewer

import (
	"fmt"

	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewForeignKeysBox(_ dtviewers.CollectionContext) *tview.Table {
	table := tview.NewTable()
	table.SetTitle("Foreign Keys")
	sneatv.DefaultBorderWithoutPadding(table.Box)
	//table.SetSelectable(true, false)

	newRow := func(row int, from, to, pk string) {
		table.SetCell(row, 0, tview.NewTableCell(from).SetTextColor(tcell.ColorLightBlue))
		table.SetCell(row, 1, tview.NewTableCell("=>").SetTextColor(tcell.ColorGray))
		table.SetCell(row, 2, tview.NewTableCell(to).SetTextColor(tcell.ColorLightGray))
		table.SetCell(row, 3, tview.NewTableCell(fmt.Sprintf("(%s)", pk)).SetTextColor(tcell.ColorSlateGray))
		return
	}
	newRow(0, "ShipVia", "Shippers", "ShipperID")
	newRow(1, "CustomerID", "Customers", "CustomerID")
	newRow(2, "EmployeeID", "Employees", "EmployeeID")
	return table
}
