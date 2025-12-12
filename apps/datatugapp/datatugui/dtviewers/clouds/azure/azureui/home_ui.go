package azureui

import (
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers/clouds"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const viewerID dtviewers.ViewerID = "azure"

func RegisterAsViewer() {
	dtviewers.RegisterViewer(dtviewers.Viewer{
		ID:          viewerID,
		Name:        "Microsoft Azure",
		Description: "(not implemented yet)",
		Shortcut:    'm',
		Action: func(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
			return GoAzureHome(&clouds.CloudContext{TUI: tui}, focusTo)
		},
	})
}

func GoAzureHome(cContext *AzureContext, focusTo sneatnav.FocusTo) error {
	menu := dtviewers.NewCloudsMenu(cContext.TUI, viewerID)

	textView := tview.NewTextView()
	sneatv.DefaultBorder(textView.Box)
	textView.SetTitle("Microsoft Azure Viewer")
	textView.SetText("Azure is not implemented yet.")

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			cContext.TUI.Header.SetFocus(sneatnav.ToBreadcrumbs, textView)
			return nil
		case tcell.KeyLeft:
			cContext.TUI.Menu.TakeFocus()
			return nil
		default:
			return event
		}
	})

	content := sneatnav.NewPanelWithBoxedPrimitive(cContext.TUI, sneatnav.WithBox(textView, textView.Box))
	cContext.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))

	return nil
}
