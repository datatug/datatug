package dbviewer

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/datatug/datatug-core/pkg/schemer"
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatcolors"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type foreignKeysBox struct {
	*tview.Table
	tui    *sneatnav.TUI
	schema schemer.ForeignKeysProvider
}

func (b foreignKeysBox) SetCollectionContext(ctx context.Context, collectionCtx dtviewers.CollectionContext) {
	b.Table.Clear()
	b.Table.SetCell(0, 0, tview.NewTableCell("Loading...").SetTextColor(tcell.ColorGray))

	go func() {
		fks, err := b.schema.GetForeignKeys(ctx, "", collectionCtx.CollectionRef.Name())
		b.tui.App.QueueUpdateDraw(func() {
			if err != nil {
				b.Table.SetCell(0, 0, tview.NewTableCell(fmt.Sprintf("Error: %v", err)).SetTextColor(tcell.ColorRed))
				return
			}
			if len(fks) == 0 {
				b.Table.SetCell(0, 0, tview.NewTableCell("No foreign keys").SetTextColor(tcell.ColorGray))
				return
			}
			for i, fk := range fks {
				b.Table.SetCell(i, 0, tview.NewTableCell(strings.Join(fk.From.Columns, ",")).SetTextColor(sneatcolors.TableColumnTitle))
				b.Table.SetCell(i, 1, tview.NewTableCell("—>"))
				b.Table.SetCell(i, 2, tview.NewTableCell(fk.To.Name))
				if !slices.Equal(fk.To.Columns, fk.From.Columns) {
					b.Table.SetCell(i, 3, tview.NewTableCell(fmt.Sprintf("(%s)", strings.Join(fk.To.Columns, ","))).SetTextColor(tview.Styles.SecondaryTextColor))
				}
			}
		})
	}()
}

func newForeignKeysBox(tui *sneatnav.TUI, schema schemer.ForeignKeysProvider) *foreignKeysBox {
	b := foreignKeysBox{
		Table:  tview.NewTable(),
		tui:    tui,
		schema: schema,
	}
	b.Table.SetTitle("Foreign Keys")
	sneatv.DefaultBorderWithoutPadding(b.Table.Box)

	return &b
}
