package dbviewer

import (
	"context"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/recordset"
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

	var rs recordset.Recordset

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC || (event.Key() == tcell.KeyRune && (event.Rune() == 'c' || event.Rune() == 'C') && (event.Modifiers()&tcell.ModMeta != 0 || event.Modifiers()&tcell.ModAlt != 0)) {
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
		case tcell.KeyUp:
			row, _ := table.GetSelection()
			if row <= 1 {
				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, table)
				return nil
			}
		case tcell.KeyLeft:
			_, col := table.GetSelection()
			if col == 0 {
				tui.SetFocus(tui.Menu)
				return nil
			}
		case tcell.KeyEnter:
			_, colIndex := table.GetSelection()
			col := rs.GetColumnByIndex(colIndex)
			name := col.Name()
			if strings.HasSuffix(name, "ID") {
				collCtx := dtviewers.CollectionContext{
					DbContext: collectionCtx.DbContext,
				}
				refTableName := name[:len(name)-len("ID")] + "s"
				collCtx.CollectionRef = dal.NewCollectionRef(refTableName, "", collectionCtx.CollectionRef.Parent())

				goTable(tui, collCtx)
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
		var err error
		rs, err = loadDataIntoTable(tui, collectionCtx, table)
		if err != nil {
			tui.App.QueueUpdateDraw(func() {
				table.SetCell(0, 0, tview.NewTableCell("Error: "+err.Error()).SetTextColor(tcell.ColorRed).SetBackgroundColor(tcell.ColorWhiteSmoke))
			})
			return
		}
	}()
}

func loadDataIntoTable(tui *sneatnav.TUI, collectionCtx dtviewers.CollectionContext, table *tview.Table) (rs recordset.Recordset, err error) {
	if collectionCtx.DbContext == nil {
		panic("collectionCtx.DbContext is nil")
	}
	db, err := collectionCtx.GetDB(context.Background())
	if err != nil {
		return nil, err
	}
	q := dal.From(collectionCtx.CollectionRef).NewQuery().SelectIntoRecordset(nil)
	ctx := context.Background()

	var tableContent TableContentRecordset
	tableContent.recordset, err = dal.ExecuteQueryAndReadAllToRecordset(ctx, q, db)
	if err != nil {
		return tableContent.recordset, err
	}
	tui.App.QueueUpdateDraw(func() {
		table.SetContent(tableContent)
		table.ScrollToBeginning()
	})
	return tableContent.recordset, nil
}
