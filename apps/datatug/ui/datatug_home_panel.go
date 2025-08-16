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

func newHomeContent(tui *sneatnav.TUI) *homePanel {
	text := tview.NewTextView()
	text.SetText("You have 2 projects.")
	panel := &homePanel{
		PanelBase: sneatnav.NewPanelBaseFromTextView(tui, text),
	}
	setPanelTitle(panel.PanelBase, "Welcome to DataTug CLI!")
	return panel
}

type homePanel struct {
	sneatnav.PanelBase
}
