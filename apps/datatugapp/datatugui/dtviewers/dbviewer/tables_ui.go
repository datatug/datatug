package dbviewer

import (
	"errors"

	"github.com/datatug/datatug-core/pkg/datatug"
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func goTables(tui *sneatnav.TUI, focusTo sneatnav.FocusTo, dbContext dtviewers.DbContext) error {
	return showCollections(tui, focusTo, dbContext, SqlDbScreenTables, "Tables", datatug.CollectionTypeTable)
}

func goViews(tui *sneatnav.TUI, focusTo sneatnav.FocusTo, dbContext dtviewers.DbContext) error {
	return showCollections(tui, focusTo, dbContext, SqlDbScreenViews, "Views", datatug.CollectionTypeView)
}

func showCollections(tui *sneatnav.TUI, focusTo sneatnav.FocusTo, dbContext dtviewers.DbContext, selectedScreen SqlDbRootScreen, title string, collectionType datatug.CollectionType) error {
	if dbContext == nil {
		return errors.New("dbContext is nil")
	}
	breadcrumbs := getSqlDbBreadcrumbs(tui, dbContext)
	breadcrumbs.Push(sneatv.NewBreadcrumb(title, nil))

	menu := newSqlDbMenu(tui, selectedScreen, dbContext)

	flex := tview.NewFlex()
	//flex.SetTitle(title + " @ " + dbContext.Driver().ShortTitle)
	//flex.SetBorder(true)

	collectionsBox := NewTablesBox(tui, dbContext, collectionType, title)
	flex.AddItem(collectionsBox, 0, 2, true)

	collectionCtx := dtviewers.CollectionContext{
		DbContext: dbContext,
	}

	columnsBox := NewColumnsBox(collectionCtx)
	flex.AddItem(columnsBox, 0, 2, true)

	flex2 := tview.NewFlex()
	flex2.SetDirection(tview.FlexRow)
	flex.AddItem(flex2, 0, 3, false)

	fks := NewForeignKeysBox(collectionCtx)
	flex2.AddItem(fks, 0, 1, false)

	referrersBox := NewReferrersBox(collectionCtx)
	flex2.AddItem(referrersBox, 0, 1, true)

	content := sneatnav.NewPanel(tui, sneatnav.WithBox(flex, flex.Box))

	tui.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))

	collectionsBox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			tui.App.SetFocus(tui.Menu)
			return nil
		case tcell.KeyRight:
			tui.App.SetFocus(columnsBox)
			return nil
		case tcell.KeyUp:
			row, _ := collectionsBox.GetSelection()
			if row <= 1 {
				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, collectionsBox)
				return nil
			}
			return event
		default:
			return event
		}
	})

	setFocusBlurFunc := func(t *tview.Table) {
		t.SetFocusFunc(func() {
			t.SetSelectable(true, false)
		})
		t.SetBlurFunc(func() {
			t.SetSelectable(false, false)
		})
	}
	setFocusBlurFunc(columnsBox)
	setFocusBlurFunc(fks)
	setFocusBlurFunc(referrersBox)

	columnsBox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			tui.App.SetFocus(fks)
			return nil
		case tcell.KeyLeft:
			tui.App.SetFocus(collectionsBox)
			return nil
		case tcell.KeyUp:
			row, _ := columnsBox.GetSelection()
			if row == 0 {
				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, columnsBox)
				return nil
			}
			return event
		default:
			return event
		}
	})

	fks.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			tui.App.SetFocus(columnsBox)
			return nil
		case tcell.KeyUp:
			row, _ := fks.GetSelection()
			if row == 0 {
				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, fks)
				return nil
			}
			return event
		case tcell.KeyDown:
			row, _ := fks.GetSelection()
			if row == fks.GetRowCount()-1 {
				tui.App.SetFocus(referrersBox)
				return nil
			}
			return event
		default:
			return event
		}
	})

	referrersBox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			tui.App.SetFocus(columnsBox)
			return nil
		case tcell.KeyUp:
			row, _ := referrersBox.GetSelection()
			if row == 0 {
				tui.App.SetFocus(fks)
				return nil
			}
			return event
		default:
			return event
		}
	})
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
