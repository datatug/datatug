package dtproject

import (
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
)

func newDashboardsPanel(ctx ProjectContext) sneatnav.Panel {
	content := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("List of dashboards here")

	sneatv.DefaultBorder(content.Box)

	return sneatnav.NewPanelWithBoxedPrimitive(ctx.TUI(), sneatnav.WithBox(content, content.Box))
}
