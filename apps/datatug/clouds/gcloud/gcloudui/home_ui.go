package gcloudui

import (
	"github.com/datatug/datatug-cli/apps/datatug/clouds"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
)

func GoHome(gcContext *GCloudContext, focusTo sneatnav.FocusTo) error {
	menu := clouds.NewCloudsMenu(gcContext.TUI, clouds.CloudGoogle)
	content := newMainMenu(gcContext, ScreenProjects)
	go func() {
		_, _ = gcContext.GetProjects()
	}()
	gcContext.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))
	return nil
}
