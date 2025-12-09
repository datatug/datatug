package dtproject

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-core/pkg/appconfig"
)

func goProjectDashboards(tui *sneatnav.TUI, project *appconfig.ProjectConfig) {
	menu := newProjectMenuPanel(tui, project, "dashboards")
	content := newDashboardsPanel(tui, project)
	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToMenu))
}
