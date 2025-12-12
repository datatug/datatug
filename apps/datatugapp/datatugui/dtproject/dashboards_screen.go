package dtproject

import (
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
)

func goProjectDashboards(tui *sneatnav.TUI, project *appconfig.ProjectConfig) {
	menu := newProjectMenuPanel(tui, project, "dashboards")
	content := newDashboardsPanel(tui, project)
	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToMenu))
}
