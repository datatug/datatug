package dtproject

import (
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
)

func goEnvironmentsScreen(tui *sneatnav.TUI, project *appconfig.ProjectConfig) {

	menu := newProjectMenuPanel(tui, project, "environments")

	content := newEnvironmentsPanel(tui, project)

	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToContent))
}
