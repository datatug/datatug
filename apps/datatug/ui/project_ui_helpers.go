package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-core/pkg/appconfig"
)

func newProjectRootScreenBase(
	tui *tapp.TUI,
	project appconfig.ProjectConfig,
	screen ProjectScreenID,
	main tapp.Panel,
	sidebar tapp.Panel,
) tapp.ScreenBase {
	grid := projectScreenGrid(tui, project, screen, main, sidebar)

	screenBase := tapp.NewScreenBase(tui, grid, tapp.FullScreen())

	screenBase.TakeFocus()

	return screenBase
}

func projectScreenGrid(
	tui *tapp.TUI,
	project appconfig.ProjectConfig,
	screenID ProjectScreenID,
	main tapp.Panel,
	sidebar tapp.Panel,
) (screen tapp.Screen) {
	_ = newProjectMenu(tui, project, screenID)

	screen, _ = newDefaultLayout(tui, projectsRootScreen, func(tui *tapp.TUI) (tapp.Panel, error) {
		return main, nil
	})

	return screen
}
