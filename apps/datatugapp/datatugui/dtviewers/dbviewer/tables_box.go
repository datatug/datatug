package dbviewer

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/datatug/datatug-core/pkg/datatug"
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewTablesBox(tui *sneatnav.TUI, dbContext dtviewers.DbContext, collectionType datatug.CollectionType, title string) *tview.Table {
	table := tview.NewTable()
	sneatv.DefaultBorderWithoutPadding(table.Box)

	// Enable cell selection by row and column
	table.SetSelectable(true, false)
	table.SetSelectedFunc(func(row, _ int) {
		cell := table.GetCell(row, 0)
		collectionInfo := cell.Reference.(*datatug.CollectionInfo)
		goTable(tui, dtviewers.CollectionContext{
			CollectionRef: collectionInfo.Ref,
			DbContext:     dbContext,
		})
	})
	{
		colIndex := 0
		addHeader := func(name string, align int, expansion int) {
			cell := tview.NewTableCell(name).
				SetSelectable(false).
				SetAlign(align).
				SetExpansion(expansion).
				SetTextColor(tcell.ColorLightBlue)
			table.SetCell(0, colIndex, cell)
			colIndex++
		}
		addHeader("Name", tview.AlignLeft, 1)
		addHeader("Records", tview.AlignRight, 0)
		table.SetFixed(1, 1)
	}

	// Start with the first data row (row 1, col 0) active
	table.SetSelectable(false, false)

	// Arrow-key behavior with edge focus transfers
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			_, col := table.GetSelection()
			if col == 0 {
				// Move focus to menu when on the leftmost column
				tui.Menu.TakeFocus()
				return nil
			}
			return event
		case tcell.KeyUp:
			row, _ := table.GetSelection()
			if row <= 1 { // row 0 is a header; row 1 is the first data row
				// Move focus to breadcrumbs when at the top row and pressing Up
				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, table)
				return nil
			}
			return event
		default:
			return event
		}
	})

	// If we create table before collections loaded columns jump

	//table.SetCell(1, 0, tview.NewTableCell("Loading...").
	//	SetSelectable(false).
	//	SetTextColor(tcell.ColorLightGrey))

	collections := make([]*datatug.CollectionInfo, 0)

	if schema := dbContext.Schema(); schema != nil {
		go func() {
			// Prime schema loading (non-blocking behavior depends on provider)
			collectionsReader, err := schema.GetCollections(context.Background(), nil)
			if err != nil {
				sneatnav.ShowErrorModal(tui, err)
				return
			}
			tui.App.QueueUpdateDraw(func() {

				i := 0

				for {
					collection, err := collectionsReader.NextCollection()
					if err != nil {
						if errors.Is(err, io.EOF) {
							break
						}
						sneatnav.ShowErrorModal(tui, err)
						return
					}
					if collection == nil {
						break
					}
					collections = append(collections, collection)
					if collection.Type() == collectionType {
						i++
						name := tview.NewTableCell(collection.Name()).SetExpansion(1)
						name.SetReference(collection)
						table.SetCell(i, 0, name)

						recordsCell := tview.NewTableCell("?").
							SetAlign(tview.AlignRight).
							SetTextColor(tcell.ColorGray)
						table.SetCell(i, 1, recordsCell)
					}
				}

				table.SetTitle(fmt.Sprintf("%s [gray](%d)", title, i))

				if i > 0 {
					table.Select(1, 0)
				}
				table.ScrollToBeginning()
				table.SetSelectable(true, false)
				table.Select(1, 0)
			})
		}()
	}

	return table
}
