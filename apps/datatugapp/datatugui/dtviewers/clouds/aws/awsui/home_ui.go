package awsui

import (
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers/clouds"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
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
	sneatv.DefaultBorderWithPadding(textView.Box)
	textView.SetTitle("Amazon Web Services")
	textView.SetText("AWS is not implemented yet.")

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			row, _ := textView.GetScrollOffset()
			if row == 0 {
				cContext.TUI.Header.SetFocus(sneatnav.ToBreadcrumbs, textView)
				return nil
			}
			return event
		case tcell.KeyLeft:
			cContext.TUI.Menu.TakeFocus()
			return nil
		default:
			return event
		}
	})

	content := sneatnav.NewPanel(cContext.TUI, sneatnav.WithBox(textView, textView.Box))
	cContext.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))

	return nil
}
