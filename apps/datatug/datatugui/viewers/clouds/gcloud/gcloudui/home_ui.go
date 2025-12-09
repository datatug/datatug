package gcloudui

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers/clouds"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
)

const viewerID viewers.ViewerID = "gc"

func RegisterAsViewer() {
	viewers.RegisterViewer(viewers.Viewer{
		ID:       viewerID,
		Name:     "Google Cloud",
		Shortcut: 'g',
		Action: func(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
			return goHome(&GCloudContext{
				CloudContext: &clouds.CloudContext{TUI: tui},
			}, focusTo)
		},
	})
}

func goHome(cContext *GCloudContext, focusTo sneatnav.FocusTo) error {
	menu := viewers.NewCloudsMenu(cContext.TUI, viewerID)
	content := newMainMenu(cContext, ScreenProjects)
	go func() {
		_, _ = cContext.GetProjects()
	}()
	cContext.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))
	return nil
}
