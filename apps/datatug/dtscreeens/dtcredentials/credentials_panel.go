package dtcredentials

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
	"github.com/datatug/datatug-cli/apps/datatug/dtnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
)

func RegisterModule() {
	datatugui.RegisterMainMenuItem(dtnav.RootScreenCredentials,
		datatugui.MainMenuItem{
			Text:     "Credentials",
			Shortcut: 'c',
			Action:   handleMenuAction,
		})
}

func handleMenuAction(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
	return goCredentials(tui, focusTo)
}

func goCredentials(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
	breadcrumbs := tui.Header.Breadcrumbs()
	breadcrumbs.Clear()
	breadcrumbs.Push(sneatv.NewBreadcrumb("Credentials", nil))
	menu := datatugui.NewDataTugMainMenu(tui, dtnav.RootScreenCredentials)
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
