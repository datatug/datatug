package datatugui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
)

func goCredentials(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
	breadcrumbs := tui.Header.Breadcrumbs()
	breadcrumbs.Clear()
	breadcrumbs.Push(sneatv.NewBreadcrumb("Credentials", nil))
	menu := newDataTugMainMenu(tui, credentialsRootScreen)
	content := newCredentialsPanel(tui)
	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToMenu))
	if focusTo == sneatnav.FocusToContent {
		tui.App.SetFocus(content)
	}
	return nil
}

func newCredentialsPanel(tui *sneatnav.TUI) sneatnav.Panel {
	text := tview.NewTextView()
	text.SetText("You have 3 credentials.")
	panel := sneatnav.NewPanelFromTextView(tui, text)
	sneatv.SetPanelTitle(panel.GetBox(), "Credentials")
	return panel
}
