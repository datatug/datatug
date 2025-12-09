package dtproject

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
)

func RegisterModule() {
	datatugui.RegisterMainMenuItem(datatugui.RootScreenProjects,
		datatugui.MainMenuItem{
			Text:     "Projects",
			Shortcut: 'p',
			Action:   GoProjectsScreen,
		})
}
