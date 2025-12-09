package dtproject

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/rivo/tview"
)

func newEnvironmentsPanel(tui *sneatnav.TUI, _ *appconfig.ProjectConfig) sneatnav.Panel {

	content := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("List of environments here")

	sneatv.DefaultBorder(content.Box)

	return sneatnav.NewPanelFromTextView(tui, content)
}
