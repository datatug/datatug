package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
)

func goCredentials(tui *sneatnav.TUI) error {
	breadcrumbs := tui.Header.Breadcrumbs()
	breadcrumbs.Clear()
	breadcrumbs.Push(sneatv.NewBreadcrumb("Credentials", nil))
	menu := newDataTugMainMenu(tui, credentialsRootScreen)
	content := newCredentialsContent(tui)
	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToMenu))
	return nil
}

func newCredentialsContent(tui *sneatnav.TUI) *credentialsPanel {
	text := tview.NewTextView()
	text.SetText("You have 3 credentials.")
	panel := &credentialsPanel{
		PanelBase: sneatnav.NewPanelBaseFromTextView(tui, text),
	}
	setPanelTitle(panel.PanelBase, "Credentials")
	return panel
}

type credentialsPanel struct {
	sneatnav.PanelBase
}
