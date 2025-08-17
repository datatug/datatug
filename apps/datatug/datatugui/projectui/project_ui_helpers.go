package projectui

//import (
//	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
//	"github.com/datatug/datatug-core/pkg/appconfig"
//)
//
//func newProjectRootScreenBase(
//	tui *sneatnav.TUI,
//	project appconfig.ProjectConfig,
//	screen ProjectScreenID,
//	main sneatnav.Panel,
//) sneatnav.ScreenBase {
//	grid := projectScreenGrid(tui, project, screen, main)
//
//	screenBase := sneatnav.NewScreenBase(tui, grid, sneatnav.FullScreen())
//
//	screenBase.TakeFocus()
//
//	return screenBase
//}
//
//func projectScreenGrid(
//	tui *sneatnav.TUI,
//	project appconfig.ProjectConfig,
//	screenID ProjectScreenID,
//	main sneatnav.Panel,
//) (screen sneatnav.Screen) {
//	_ = NewProjectMenuPanel(tui, project, screenID)
//
//	screen = newDefaultLayout(tui, projectsRootScreen, func(tui *sneatnav.TUI) (sneatnav.Panel, error) {
//		return main, nil
//	})
//
//	return screen
//}
