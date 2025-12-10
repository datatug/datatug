package sqlviewer

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SqlDbRootScreen int

const (
	SqlDbScreenTables SqlDbRootScreen = iota
	SqlDbScreenViews
)

func newSqlDbMenu(tui *sneatnav.TUI, selectedScreen SqlDbRootScreen, filePath string) sneatnav.Panel {
	list := tview.NewList()
	list.SetWrapAround(false)

	list.AddItem("Tables", "", 't', func() {
		_ = goTables(tui, sneatnav.FocusToContent, filePath)
	})
	list.AddItem("Views", "", 'v', func() {
		_ = goViews(tui, sneatnav.FocusToContent, filePath)
	})

	list.SetCurrentItem(int(selectedScreen))

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			tui.Content.TakeFocus()
			return nil
		case tcell.KeyUp:
			if list.GetCurrentItem() == 0 {
				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, list)
				return nil
			}
		default:
			return event
		}
		return event
	})

	return sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(list, list.Box))
}
