package awsui

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers/clouds"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const viewerID viewers.ViewerID = "aws"

func RegisterAsViewer() {
	viewers.RegisterViewer(viewers.Viewer{
		ID:       viewerID,
		Name:     "Amazon Web Services",
		Shortcut: 'a',
		Action: func(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
			return goAwsHome(&clouds.CloudContext{TUI: tui}, focusTo)
		},
	})
}

func goAwsHome(cContext *AwsContext, focusTo sneatnav.FocusTo) error {
	menu := viewers.NewCloudsMenu(cContext.TUI, viewerID)

	textView := tview.NewTextView()
	sneatv.DefaultBorder(textView.Box)
	textView.SetTitle("Amazon Web Services")
	textView.SetText("AWS is not implemented yet.")

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
