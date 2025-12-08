package awsui

import (
	"github.com/datatug/datatug-cli/apps/datatug/clouds"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GoAwsHome(cContext *AwsContext, focusTo sneatnav.FocusTo) error {
	menu := clouds.NewCloudsMenu(cContext.TUI, clouds.CloudAWS)

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
