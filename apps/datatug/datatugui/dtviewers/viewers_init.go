package dtviewers

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
	"github.com/datatug/datatug-cli/apps/datatug/dtnav"
)

var viewers []Viewer

func RegisterViewer(viewer Viewer) {
	viewers = append(viewers, viewer)
}

func RegisterModule() {
	datatugui.RegisterMainMenuItem(dtnav.RootScreenViewers,
		datatugui.MainMenuItem{
			Text:     "Viewers",
			Shortcut: 'v',
			Action:   goViewersScreen,
		})
}
