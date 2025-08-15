package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-cli/pkg/tvprimitives/breadcrumbs"
	"github.com/rivo/tview"
)

func goCredentials(tui *tapp.TUI) error {
	tui.Header.Breadcrumbs.Clear()
	tui.Header.Breadcrumbs.Push(breadcrumbs.NewBreadcrumb("Credentials", nil))
	menu := newDataTugMainMenu(tui, credentialsRootScreen)
	content := newCredentialsContent(tui)
	tui.SetPanels(menu, content)
	return nil
}

func newCredentialsContent(tui *tapp.TUI) *credentialsPanel {
	text := tview.NewTextView()
	text.SetText("You have 3 credentials.")
	panel := &credentialsPanel{
		PanelBase: tapp.NewPanelBaseFromTextView(tui, text),
	}
	setPanelTitle(panel.PanelBase, "Credentials")
	return panel
}

type credentialsPanel struct {
	tapp.PanelBase
}
