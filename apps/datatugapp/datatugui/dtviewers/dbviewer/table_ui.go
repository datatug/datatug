package dbviewer

import (
	"context"
	"fmt"
	"reflect"

	"github.com/dal-go/dalgo/dal"
	"github.com/datatug/datatug-core/pkg/schemer"
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/rivo/tview"
)

func goTable(tui *sneatnav.TUI, collectionCtx dtviewers.CollectionContext) {
	table := tview.NewTable()

	table.SetTitle(collectionCtx.CollectionRef.Name())

	content := sneatnav.NewPanel(tui, sneatnav.WithBox(table, table.Box))

	tui.SetPanels(nil, content)

	go func() {
		_ = loadDataIntoTable(collectionCtx, table)
	}()
}

func loadDataIntoTable(collectionCtx dtviewers.CollectionContext, table *tview.Table) (err error) {
	db, err := collectionCtx.GetDB(context.Background())
	if err != nil {
		return err
	}
	q := dal.NewQueryBuilder(dal.From(collectionCtx.CollectionRef)).SelectInto(func() dal.Record {
		return dal.NewRecordWithIncompleteKey(collectionCtx.CollectionRef.Name(), reflect.Invalid, make(map[string]any))
	})
	ctx := context.Background()

	records, err := q.ReadRecords(ctx, db)
	if err != nil {
		return err
	}

	schema := collectionCtx.Schema()
	// TODO: Pass CollectionRef to GetColumns() by value?
	columns, err := schema.GetColumns(ctx, "", schemer.ColumnsFilter{CollectionRef: &collectionCtx.CollectionRef})
	if err != nil {
		return err
	}

	for i, record := range records {
		data := record.Data().(map[string]any)
		for j, col := range columns {
			v := data[col.Name]
			table.SetCell(i+1, j, tview.NewTableCell(fmt.Sprintf("%v", v)))
		}
	}

	return nil
}
