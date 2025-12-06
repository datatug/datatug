package dtviewers

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
	"github.com/datatug/datatug-cli/apps/datatug/dtnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func goViewersScreen(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
	breadcrumbs := tui.Header.Breadcrumbs()
	breadcrumbs.Clear()
	breadcrumbs.Push(sneatv.NewBreadcrumb("Viewers", nil))

	// Set secondary text color to gray
	viewersList.SetSecondaryTextColor(tcell.ColorDarkGray)

	viewersList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC, tcell.KeyBacktab, tcell.KeyLeft:
			tui.SetFocus(tui.Menu)
			return nil
		case tcell.KeyUp:
			if viewersList.GetCurrentItem() == 0 {
				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, viewersList)
				return nil
			}
			return event
		case tcell.KeyDown:
			// Prevent jumping to first item when on last item
			if viewersList.GetCurrentItem() == viewersList.GetItemCount()-1 {
				return nil
			}
			return event
		default:
			return event
		}
	})

	sneatv.DefaultBorder(viewersList.Box)
	// Set spacing between items to 1 line by increasing vertical padding
	viewersList.SetBorderPadding(1, 1, 1, 1)
	viewersList.SetTitle(" Viewers ")
	viewersList.SetTitleAlign(tview.AlignLeft)

	menu := datatugui.NewDataTugMainMenu(tui, dtnav.RootScreenViewers)
	content := sneatnav.NewPanelFromList(tui, viewersList)

	tui.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))
	return nil
}
