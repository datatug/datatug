package dbviewer

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func goTables(tui *sneatnav.TUI, _ sneatnav.FocusTo, dbContext dtviewers.DbContext) error {
	return showCollections(tui, dbContext, SqlDbScreenTables, "Tables", "table")
}

func goViews(tui *sneatnav.TUI, _ sneatnav.FocusTo, dbContext dtviewers.DbContext) error {
	return showCollections(tui, dbContext, SqlDbScreenViews, "Views", "view")
}

func showCollections(tui *sneatnav.TUI, dbContext dtviewers.DbContext, selectedScreen SqlDbRootScreen, title, collectionType string) error {
	if dbContext == nil {
		return errors.New("dbContext is nil")
	}
	breadcrumbs := getSqlDbBreadcrumbs(tui, dbContext)
	breadcrumbs.Push(sneatv.NewBreadcrumb(title, nil))

	menu := newSqlDbMenu(tui, selectedScreen, dbContext)

	table := tview.NewTable()
	table.SetTitle(title + " @ " + dbContext.Driver().ShortTitle)
	// Enable cell selection by row and column
	table.SetSelectable(true, true)
	// Start with the first data row (row 1, col 0) active
	table.Select(1, 0)
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
			if row <= 1 { // row 0 is header; row 1 is first data row
				// Move focus to breadcrumbs when at the top row and pressing Up
				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, table)
				return nil
			}
			return event
		default:
			return event
		}
	})

	colIndex := 0
	addHeader := func(name string) {
		cell := tview.NewTableCell(name).
			SetSelectable(false).
			SetTextColor(tcell.ColorLightBlue)
		table.SetCell(0, colIndex, cell)
		colIndex++
	}
	addHeader("Name")
	addHeader("Cols")

	table.SetCell(1, 0, tview.NewTableCell("Loading...").
		SetSelectable(false).
		SetTextColor(tcell.ColorLightGrey))

	go func() {
		if schema := dbContext.Schema(); schema != nil {
			// Prime schema loading (non-blocking behavior depends on provider)
			collections, err := schema.GetCollections(context.Background(), nil)
			if err != nil {
				table.Clear()
				sneatnav.ShowErrorModal(tui, err)
				return
			}
			tui.App.QueueUpdateDraw(func() {
				count := 0
				i := 0

				for {
					collection, err := collections.NextCollection()
					if err != nil {
						sneatnav.ShowErrorModal(tui, err)
						return
					}
					if collection == nil {
						break
					}
					count++
					if collection.DbType == collectionType {
						i++
						name := tview.NewTableCell(collection.Name)
						table.SetCell(i, 0, name)

						cols := tview.NewTableCell(strconv.Itoa(i)).SetAlign(tview.AlignRight)
						table.SetCell(i, 1, cols)
					}
				}
				table.SetTitle(fmt.Sprintf("%d %s @ %s", count, title, dbContext.Driver().ShortTitle))

				if count > 0 {
					table.Select(1, 0)
				}
				table.ScrollToBeginning()
			})
		}
	}()

	content := sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(table, table.Box))

	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToContent))
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
