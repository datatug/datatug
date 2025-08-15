package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-cli/pkg/tvprimitives/breadcrumbs"
	"github.com/rivo/tview"
)

func GoHomeScreen(tui *tapp.TUI) error {
	tui.Header.Breadcrumbs.Clear()
	tui.Header.Breadcrumbs.Push(breadcrumbs.NewBreadcrumb("Home", nil))
	menu := newDataTugMainMenu(tui, homeRootScreen)
	content := newHomeContent(tui)
	tui.SetPanels(menu, content)
	return nil
}

func newHomeContent(tui *tapp.TUI) *homePanel {
	text := tview.NewTextView()
	text.SetText("You have 2 projects.")
	panel := &homePanel{
		PanelBase: tapp.NewPanelBaseFromTextView(tui, text),
	}
	setPanelTitle(panel.PanelBase, "Welcome to DataTug CLI!")
	return panel
}

type homePanel struct {
	tapp.PanelBase
}
