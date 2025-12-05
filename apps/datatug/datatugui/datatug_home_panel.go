package datatugui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GoHomeScreen(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
	breadcrumbs := tui.Header.Breadcrumbs()
	breadcrumbs.Clear()
	breadcrumbs.Push(sneatv.NewBreadcrumb("Home", nil))
	menu := newDataTugMainMenu(tui, homeRootScreen)
	content := newHomeContent(tui)
	tui.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))
	return nil
}

func newHomeContent(tui *sneatnav.TUI) sneatnav.Panel {
	text := tview.NewTextView()
	text.SetText("You have 2 projects.")
	sneatv.SetPanelTitle(text.Box, "Welcome to DataTug CLI!")
	text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC, tcell.KeyBacktab, tcell.KeyLeft:
			tui.SetFocus(tui.Menu)
			return nil
		default:
			return event
		}
	})
	return sneatnav.NewPanelFromTextView(tui, text)
}
