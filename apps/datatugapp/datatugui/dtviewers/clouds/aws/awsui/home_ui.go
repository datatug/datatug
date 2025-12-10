package awsui

import (
	"github.com/datatug/datatug-cli/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug-cli/apps/datatugapp/datatugui/dtviewers/clouds"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const viewerID dtviewers.ViewerID = "aws"

func RegisterAsViewer() {
	dtviewers.RegisterViewer(dtviewers.Viewer{
		ID:          viewerID,
		Name:        "Amazon Web Services",
		Description: "(not implemented yet)",
		Shortcut:    'a',
		Action: func(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
			return goAwsHome(&clouds.CloudContext{TUI: tui}, focusTo)
		},
	})
}

func goAwsHome(cContext *AwsContext, focusTo sneatnav.FocusTo) error {
	menu := dtviewers.NewCloudsMenu(cContext.TUI, viewerID)

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

	content := sneatnav.NewPanelWithBoxedPrimitive(cContext.TUI, sneatnav.WithBox(textView, textView.Box))
	cContext.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))

	return nil
}
