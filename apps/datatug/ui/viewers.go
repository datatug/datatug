package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var _ tview.Primitive = (*viewersPanel)(nil)
var _ sneatnav.Cell = (*viewersPanel)(nil)

type viewersPanel struct {
	sneatnav.PanelBase
	list *tview.List
}

func (p *viewersPanel) Draw(screen tcell.Screen) {
	p.list.Draw(screen)
}

func goViewersScreen(tui *sneatnav.TUI) error {
	breadcrumbs := tui.Header.Breadcrumbs()
	breadcrumbs.Clear()
	breadcrumbs.Push(sneatv.NewBreadcrumb("Viewers", nil))
	list := tview.NewList()

	// Add the two required items
	list.AddItem("Firestore viewer", "Browse & edit data in Firestore databases", '1', nil)
	list.AddItem("SQL DB viewer", "Browse & query SQL databases", '2', nil)

	// Set secondary text color to gray
	list.SetSecondaryTextColor(tcell.ColorDarkGray)

	content := &viewersPanel{
		PanelBase: sneatnav.NewPanelBaseFromList(tui, list),
		list:      list,
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC, tcell.KeyBacktab, tcell.KeyLeft:
			tui.SetFocus(tui.Menu)
			return nil
		case tcell.KeyUp:
			if list.GetCurrentItem() == 0 {

				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, list)
				return nil
			}
			return event
		case tcell.KeyDown:
			// Prevent jumping to first item when on last item
			if list.GetCurrentItem() == list.GetItemCount()-1 {
				return nil
			}
			return event
		default:
			return event
		}
	})

	defaultBorder(content.list.Box)
	// Set spacing between items to 1 line by increasing vertical padding
	content.list.SetBorderPadding(1, 1, 1, 1)
	content.list.SetTitle(" Viewers ")
	content.list.SetTitleAlign(tview.AlignLeft)

	menu := newDataTugMainMenu(tui, viewersRootScreen)
	tui.SetPanels(menu, content)
	return nil
}
