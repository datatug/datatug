package dthome

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
	"github.com/datatug/datatug-cli/apps/datatug/dtnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RegisterModule() {
	datatugui.RegisterMainMenuItem(dtnav.RootScreenHome,
		datatugui.MainMenuItem{
			Text:     "Home",
			Shortcut: 'h',
			Action:   GoHomeScreen,
		})
}

func GoHomeScreen(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
	breadcrumbs := tui.Header.Breadcrumbs()
	breadcrumbs.Clear()
	breadcrumbs.Push(sneatv.NewBreadcrumb("Home", nil))
	menu := datatugui.NewDataTugMainMenu(tui, dtnav.RootScreenHome)
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
		case tcell.KeyDown:
			tui.SetFocus(tui.Menu)
			// TODO(help-wanted): Ideally we'd want to move to next main menu item but this does not happen
			return event
		default:
			return event
		}
	})
	return sneatnav.NewPanelFromTextView(tui, text)
}
