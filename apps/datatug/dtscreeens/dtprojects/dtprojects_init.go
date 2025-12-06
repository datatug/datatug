package dtprojects

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
	"github.com/datatug/datatug-cli/apps/datatug/dtnav"
)

func RegisterModule() {
	datatugui.RegisterMainMenuItem(dtnav.RootScreenProjects,
		datatugui.MainMenuItem{
			Text:     "Projects",
			Shortcut: 'p',
			Action:   GoProjectsScreen,
		})
}
