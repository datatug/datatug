package dbviewer

import (
	"context"

	"github.com/atotto/clipboard"
	"github.com/dal-go/dalgo/dal"
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func goTable(tui *sneatnav.TUI, collectionCtx dtviewers.CollectionContext) {

	tableName := collectionCtx.CollectionRef.Name()
	breadcrumbs := getSqlDbBreadcrumbs(tui, collectionCtx.DbContext)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Tables", nil))
	breadcrumbs.Push(sneatv.NewBreadcrumb(tableName, nil))

	table := tview.NewTable()
	table.SetTitle("Table: " + tableName)
	table.SetSelectable(true, true)
	table.SetFixed(1, 0)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC || (event.Key() == tcell.KeyRune && event.Rune() == 'c' && event.Modifiers()&tcell.ModMeta != 0) {
			row, col := table.GetSelection()
			if row >= 0 && col >= 0 {
				cell := table.GetCell(row, col)
				if cell != nil && cell.Text != "" {
					_ = clipboard.WriteAll(cell.Text)
				}
			}
			return nil
		}

		switch event.Key() {
		case tcell.KeyLeft:
			_, col := table.GetSelection()
			if col == 0 {
				tui.SetFocus(tui.Menu)
				return nil
			}
		default:
			return event
		}
		return event
	})

	menu := newSqlDbMenu(tui, SqlDbScreenTables, collectionCtx.DbContext)
	content := sneatnav.NewPanel(tui, sneatnav.WithBox(table, table.Box))

	tui.SetPanels(menu, content)

	go func() {
		err := loadDataIntoTable(tui, collectionCtx, table)
		if err != nil {
			table.SetCell(0, 0, tview.NewTableCell("Error: "+err.Error()).SetTextColor(tcell.ColorRed).SetBackgroundColor(tcell.ColorWhiteSmoke))
			return
		}
	}()
}

func loadDataIntoTable(tui *sneatnav.TUI, collectionCtx dtviewers.CollectionContext, table *tview.Table) (err error) {
	if collectionCtx.DbContext == nil {
		panic("collectionCtx.DbContext is nil")
	}
	db, err := collectionCtx.GetDB(context.Background())
	if err != nil {
		return err
	}
	q := dal.From(collectionCtx.CollectionRef).NewQuery().SelectIntoRecordset(nil)
	ctx := context.Background()

	var tableContent TableContentRecordset
	tableContent.recordset, err = dal.ExecuteQueryAndReadAllToRecordset(ctx, q, db)
	if err != nil {
		return err
	}
	tui.App.QueueUpdateDraw(func() {
		table.SetContent(tableContent)
	})
	return nil
}
