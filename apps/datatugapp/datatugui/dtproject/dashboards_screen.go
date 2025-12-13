package dtproject

import (
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
)

func goProjectDashboards(ctx ProjectContext) {
	menu := newProjectMenuPanel(ctx, "dashboards")
	content := newDashboardsPanel(ctx)
	ctx.TUI().SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToMenu))
}
