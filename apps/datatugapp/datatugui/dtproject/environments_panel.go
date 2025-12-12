package dtproject

import (
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/rivo/tview"
)

func newEnvironmentsPanel(tui *sneatnav.TUI, _ *appconfig.ProjectConfig) sneatnav.Panel {
	textView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	textView.SetText("List of environments here")
	return sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(textView, textView.Box))
}
