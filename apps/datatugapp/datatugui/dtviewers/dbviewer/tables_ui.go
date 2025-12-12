package dbviewer

import (
	"context"
	"fmt"

	"github.com/datatug/datatug-cli/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
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
	breadcrumbs := getSqlDbBreadcrumbs(tui, dbContext)
	breadcrumbs.Push(sneatv.NewBreadcrumb(title, nil))

	menu := newSqlDbMenu(tui, selectedScreen, dbContext)

	list := tview.NewList()
	list.SetTitle(title + " @ " + dbContext.Driver().ShortTitle)
	list.SetWrapAround(false)
	setDefaultInputCaptureForList(tui, list)

	list.AddItem("Loading...", "Please wait.", 0, nil)

	if schema := dbContext.Schema(); schema != nil {
		go func() {
			// Prime schema loading (non-blocking behavior depends on provider)
			collections, err := schema.GetCollections(context.Background(), nil)
			if err != nil {
				list.Clear()
				list.AddItem("Error", err.Error(), 0, nil)
				return
			}
			tui.App.QueueUpdateDraw(func() {
				list.Clear()
				count := 0
				for {
					collection, err := collections.NextCollection()
					if err != nil {
						list.AddItem("Error", err.Error(), 0, nil)
						return
					}
					if collection == nil {
						break
					}
					count++
					if collection.DbType == collectionType {
						list.AddItem(collection.Name, "", 0, nil)
					}
				}
				list.SetTitle(fmt.Sprintf("%d %s @ %s", count, title, dbContext.Driver().ShortTitle))
			})
		}()
	}

	content := sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(list, list.Box))

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
