package gcloudui

import (
	"github.com/datatug/datatug-cli/apps/datatug/clouds"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
)

func GoHome(cContext *GCloudContext, focusTo sneatnav.FocusTo) error {
	menu := clouds.NewCloudsMenu(cContext.TUI, clouds.CloudGoogle)
	content := newMainMenu(cContext, ScreenProjects)
	go func() {
		_, _ = cContext.GetProjects()
	}()
	cContext.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))
	return nil
}
