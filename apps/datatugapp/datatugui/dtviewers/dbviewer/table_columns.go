package dbviewer

import (
	"context"

	"github.com/datatug/datatug-core/pkg/schemer"
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewColumnsBox(ctx context.Context, collectionCtx dtviewers.CollectionContext, tui *sneatnav.TUI) *tview.Table {

	table := tview.NewTable()
	table.SetTitle(`[gray]Categories:[-] Columns`)
	table.SetFixed(1, 1)
	sneatv.DefaultBorderWithoutPadding(table.Box)

	addHeader := func() {
		table.SetCell(0, 0, tview.NewTableCell("Name").SetExpansion(1))
		table.SetCell(0, 1, tview.NewTableCell("Type"))
	}
	addHeader()
	table.SetCell(1, 0, tview.NewTableCell("Loading...").SetTextColor(tcell.ColorGray))

	schema := collectionCtx.Schema()
	if schema == nil {
		return nil
	}

	go func() {
		columns, err := schema.GetColumns(ctx, "", schemer.ColumnsFilter{
			CollectionRef: &collectionCtx.CollectionRef,
		})

		tui.App.QueueUpdateDraw(func() {
			table.Clear()
			if err != nil {
				table.SetCell(0, 0, tview.NewTableCell(err.Error()).SetTextColor(tcell.ColorRed))
				return
			}
			addHeader()
			for i, col := range columns {
				table.SetCell(i+1, 0, tview.NewTableCell(col.Name))
				table.SetCell(i+1, 1, tview.NewTableCell(col.DbType).SetTextColor(tcell.ColorGray))
			}
		})
	}()

	return table
}
