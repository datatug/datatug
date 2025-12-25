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

type columnsBox struct {
	collectionCtx dtviewers.CollectionContext
	schema        schemer.ColumnsProvider
	tui           *sneatnav.TUI
	*tview.Table
}

func (b columnsBox) SetCollectionContext(ctx context.Context, collectionCtx dtviewers.CollectionContext) {
	b.Clear()
	b.collectionCtx = collectionCtx
	b.addHeader()
	b.SetCell(1, 0, tview.NewTableCell("Loading...").SetTextColor(tcell.ColorGray))

	go func() {
		columns, err := b.schema.GetColumns(ctx, "", schemer.ColumnsFilter{
			CollectionRef: &collectionCtx.CollectionRef,
		})

		b.tui.App.QueueUpdateDraw(func() {
			b.Table.Clear()
			if err != nil {
				b.Table.SetCell(0, 0, tview.NewTableCell(err.Error()).SetTextColor(tcell.ColorRed))
				return
			}
			b.addHeader()
			for i, col := range columns {
				b.Table.SetCell(i+1, 0, tview.NewTableCell(col.Name))
				b.Table.SetCell(i+1, 1, tview.NewTableCell(col.DbType).SetTextColor(tcell.ColorGray))
			}
		})
	}()
}

func (b columnsBox) addHeader() {
	b.SetCell(0, 0, tview.NewTableCell("Name").SetExpansion(1))
	b.SetCell(0, 1, tview.NewTableCell("Type"))
}

func newColumnsBox(ctx context.Context, dbContext dtviewers.DbContext, tui *sneatnav.TUI) (b *columnsBox) {
	schema := dbContext.Schema()
	if schema == nil {
		return nil
	}

	b = &columnsBox{
		schema: schema,
		tui:    tui,
		Table:  tview.NewTable().SetFixed(1, 1),
	}
	b.Table.SetTitle(`[gray]Categories:[-] Columns`)
	sneatv.DefaultBorderWithoutPadding(b.Table.Box)

	return
}
