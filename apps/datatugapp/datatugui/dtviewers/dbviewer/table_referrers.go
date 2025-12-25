package dbviewer

import (
	"context"
	"fmt"
	"strings"

	"github.com/datatug/datatug-core/pkg/schemer"
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatcolors"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type referrersBox struct {
	tui *sneatnav.TUI
	*tview.Table
	schema schemer.ReferrersProvider
}

func (b *referrersBox) SetCollectionContext(ctx context.Context, collectionCtx dtviewers.CollectionContext) {
	b.Table.Clear()
	b.Table.SetCell(0, 0, tview.NewTableCell("Loading...").SetTextColor(tcell.ColorGray))

	go func() {
		referrers, err := b.schema.GetReferrers(ctx, "", collectionCtx.CollectionRef.Name())
		b.tui.App.QueueUpdateDraw(func() {
			if err != nil {
				b.Table.SetCell(0, 0, tview.NewTableCell(fmt.Sprintf("Error: %v", err)).SetTextColor(tcell.ColorRed))
				return
			}
			if len(referrers) == 0 {
				b.Table.SetCell(0, 0, tview.NewTableCell("No referrers").SetTextColor(tcell.ColorGray))
				return
			}
			for i, referrer := range referrers {
				b.Table.SetCell(i, 0, tview.NewTableCell("<=").SetTextColor(tcell.ColorGray))
				b.Table.SetCell(i, 1, tview.NewTableCell(referrer.From.Name).SetTextColor(sneatcolors.TableColumnTitle))
				b.Table.SetCell(i, 2, tview.NewTableCell(fmt.Sprintf("(%s)", strings.Join(referrer.From.Columns, ","))))
			}
		})
	}()
}

func newReferrersBox(tui *sneatnav.TUI, schema schemer.ReferrersProvider) *referrersBox {
	b := referrersBox{
		tui:    tui,
		schema: schema,
		Table:  tview.NewTable(),
	}
	b.Table.SetTitle("Referrers")
	sneatv.DefaultBorderWithoutPadding(b.Table.Box)

	return &b
}
