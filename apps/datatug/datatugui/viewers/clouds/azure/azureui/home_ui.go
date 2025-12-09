package azureui

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers/clouds"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const viewerID viewers.ViewerID = "azure"

func RegisterAsViewer() {
	viewers.RegisterViewer(viewers.Viewer{
		ID:       viewerID,
		Name:     "Microsoft Azure",
		Shortcut: 'm',
		Action: func(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
			return GoAzureHome(&clouds.CloudContext{TUI: tui}, focusTo)
		},
	})
}

func GoAzureHome(cContext *AzureContext, focusTo sneatnav.FocusTo) error {
	menu := viewers.NewCloudsMenu(cContext.TUI, viewerID)

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

	content := sneatnav.NewPanelFromTextView(cContext.TUI, textView)
	cContext.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))

	return nil
}
