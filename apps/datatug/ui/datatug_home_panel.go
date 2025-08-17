package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
)

func GoHomeScreen(tui *sneatnav.TUI) error {
	breadcrumbs := tui.Header.Breadcrumbs()
	breadcrumbs.Clear()
	breadcrumbs.Push(sneatv.NewBreadcrumb("Home", nil))
	menu := newDataTugMainMenu(tui, homeRootScreen)
	content := newHomeContent(tui)
	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToMenu))
	return nil
}

func newHomeContent(tui *sneatnav.TUI) sneatnav.Panel {
	text := tview.NewTextView()
	text.SetText("You have 2 projects.")
	setPanelTitle(text.Box, "Welcome to DataTug CLI!")
	return sneatnav.NewPanelFromTextView(tui, text)
}
