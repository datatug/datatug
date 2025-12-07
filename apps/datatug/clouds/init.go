package clouds

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
	"github.com/datatug/datatug-cli/apps/datatug/dtnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
)

type Cloud struct {
	Name     string
	Shortcut rune
	Action   func(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error
}

var registeredClouds []Cloud

func RegisterModule(clouds []Cloud) {
	registeredClouds = clouds
	datatugui.RegisterMainMenuItem(dtnav.RootScreenClouds,
		datatugui.MainMenuItem{
			Text:     "Clouds",
			Shortcut: 'c',
			Action:   goClouds,
		})
}
