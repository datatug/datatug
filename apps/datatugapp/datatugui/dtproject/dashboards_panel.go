package dtproject

import (
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
)

func newDashboardsPanel(tui *sneatnav.TUI, _ *appconfig.ProjectConfig) sneatnav.Panel {
	content := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("List of dashboards here")

	sneatv.DefaultBorder(content.Box)

	return sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(content, content.Box))
}
