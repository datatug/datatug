package dbviewer

import (
	"fmt"

	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewForeignKeysBox(_ dtviewers.CollectionContext) *tview.Table {
	table := tview.NewTable()
	table.SetTitle("Foreign Keys")
	table.SetBorder(true)
	table.SetBorderColor(tcell.ColorDarkSlateGray)
	//table.SetSelectable(true, false)

	newRow := func(row int, from, to, pk string) {
		table.SetCell(row, 0, tview.NewTableCell(from))
		table.SetCell(row, 1, tview.NewTableCell("=>"))
		table.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%s(%s)", to, pk)))
		return
	}
	newRow(0, "ShipVia", "Shippers", "ShipperID")
	newRow(1, "CustomerID", "Customers", "CustomerID")
	newRow(2, "EmployeeID", "Employees", "EmployeeID")
	return table
}
