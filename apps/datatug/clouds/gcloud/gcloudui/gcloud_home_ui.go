package gcloudui

import (
	"github.com/datatug/datatug-cli/apps/datatug/clouds"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
)

func GoHome(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
	menu := clouds.NewCloudsMenu(tui, clouds.CloudGoogle)
	content := newMainMenu(tui, ScreenProjects)

	tui.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))
	return nil
}
