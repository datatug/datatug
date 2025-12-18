package dbviewer

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/datatug/datatug-core/pkg/datatug"
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func goTables(tui *sneatnav.TUI, focusTo sneatnav.FocusTo, dbContext dtviewers.DbContext) error {
	return showCollections(tui, focusTo, dbContext, SqlDbScreenTables, "Tables", "table")
}

func goViews(tui *sneatnav.TUI, focusTo sneatnav.FocusTo, dbContext dtviewers.DbContext) error {
	return showCollections(tui, focusTo, dbContext, SqlDbScreenViews, "Views", "view")
}

func showCollections(tui *sneatnav.TUI, focusTo sneatnav.FocusTo, dbContext dtviewers.DbContext, selectedScreen SqlDbRootScreen, title, collectionType string) error {
	if dbContext == nil {
		return errors.New("dbContext is nil")
	}
	breadcrumbs := getSqlDbBreadcrumbs(tui, dbContext)
	breadcrumbs.Push(sneatv.NewBreadcrumb(title, nil))

	menu := newSqlDbMenu(tui, selectedScreen, dbContext)

	collectionsTable := tview.NewTable()
	collectionsTable.SetTitle(title + " @ " + dbContext.Driver().ShortTitle)
	// Enable cell selection by row and column
	collectionsTable.SetSelectable(true, true)
	// Start with the first data row (row 1, col 0) active
	collectionsTable.Select(1, 0)
	// Arrow-key behavior with edge focus transfers
	collectionsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			_, col := collectionsTable.GetSelection()
			if col == 0 {
				// Move focus to menu when on the leftmost column
				tui.Menu.TakeFocus()
				return nil
			}
			return event
		case tcell.KeyUp:
			row, _ := collectionsTable.GetSelection()
			if row <= 1 { // row 0 is header; row 1 is first data row
				// Move focus to breadcrumbs when at the top row and pressing Up
				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, collectionsTable)
				return nil
			}
			return event
		default:
			return event
		}
	})

	// If we create collectionsTable before collections loaded columns jump

	//collectionsTable.SetCell(1, 0, tview.NewTableCell("Loading...").
	//	SetSelectable(false).
	//	SetTextColor(tcell.ColorLightGrey))

	collections := make([]*datatug.CollectionInfo, 0)

	go func() {
		if schema := dbContext.Schema(); schema != nil {
			// Prime schema loading (non-blocking behavior depends on provider)
			collectionsReader, err := schema.GetCollections(context.Background(), nil)
			if err != nil {
				sneatnav.ShowErrorModal(tui, err)
				return
			}
			tui.App.QueueUpdateDraw(func() {
				colIndex := 0
				addHeader := func(name string) {
					cell := tview.NewTableCell(name).
						SetSelectable(false).
						SetTextColor(tcell.ColorLightBlue)
					collectionsTable.SetCell(0, colIndex, cell)
					colIndex++
				}
				addHeader("Name")
				addHeader("Cols")
				collectionsTable.SetFixed(1, 1)

				i := 0

				for {
					collection, err := collectionsReader.NextCollection()
					collections = append(collections, collection)
					if err != nil {
						sneatnav.ShowErrorModal(tui, err)
						return
					}
					if collection == nil {
						break
					}
					if collection.DbType == collectionType {
						i++
						name := tview.NewTableCell(collection.Name)
						name.SetReference(collection)
						collectionsTable.SetCell(i, 0, name)

						cols := tview.NewTableCell(strconv.Itoa(i)).SetAlign(tview.AlignRight)
						collectionsTable.SetCell(i, 1, cols)
					}
				}
				collectionsTable.SetTitle(fmt.Sprintf("%d %s @ %s", i, title, dbContext.Driver().ShortTitle))

				if i > 0 {
					collectionsTable.Select(1, 0)
				}
				collectionsTable.ScrollToBeginning()
			})
		}
	}()

	//sidePanel := tview.NewFlex().SetDirection(tview.FlexRow)
	//
	//primaryKey := tview.NewTextView()
	////primaryKey.SetBorder(true)
	//primaryKey.SetText("Primary Key: Loading...")
	//
	//columns := tview.NewTextView()
	//columns.SetBorder(true)
	//columns.SetTitle("Columns")
	//columns.SetTextColor(tcell.ColorLightGrey)
	//columns.SetText("Loading...")
	//
	//sidePanel.AddItem(primaryKey, 1, 0, false)
	//sidePanel.AddItem(columns, 0, 1, false)
	//
	//contentFlex := tview.NewFlex()
	//
	//contentFlex.AddItem(collectionsTable, 0, 1, true)
	//contentFlex.AddItem(sidePanel, 0, 1, true)

	content := sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(collectionsTable, collectionsTable.Box))

	tui.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))
	return nil
}

func setDefaultInputCaptureForList(tui *sneatnav.TUI, list *tview.List) {
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			tui.Menu.TakeFocus()
			return nil
		case tcell.KeyUp:
			if list.GetCurrentItem() == 0 {
				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, list)
				return nil
			}
			return event
		default:
			return event
		}
	})
}

//func setDefaultInputCapture(tui *sneatnav.TUI, c interface {
//	tview.Primitive
//	SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *tview.Box
//}) {
//	c.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
//		switch event.Key() {
//		case tcell.KeyLeft:
//			tui.Menu.TakeFocus()
//			return nil
//		case tcell.KeyUp:
//			tui.Header.SetFocus(sneatnav.ToBreadcrumbs, c)
//			return nil
//		default:
//			return event
//		}
//	})
//}
