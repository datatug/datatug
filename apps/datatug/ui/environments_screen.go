package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-core/pkg/appconfig"
)

func goEnvironmentsScreen(tui *sneatnav.TUI, project *appconfig.ProjectConfig) {

	menu := newProjectsMenuPanel(tui)
	content := newEnvironmentsPanel(tui, project)
	tui.SetPanels(menu, content)
}
