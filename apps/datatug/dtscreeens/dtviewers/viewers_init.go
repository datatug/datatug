package dtviewers

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
	"github.com/datatug/datatug-cli/apps/datatug/dtnav"
	"github.com/rivo/tview"
)

func init() {
	AddViewer("SQL DB viewer", "Browse & query SQL databases", '2', nil)
}

var viewersList = tview.NewList()

var AddViewer = viewersList.AddItem

func RegisterModule() {
	datatugui.RegisterMainMenuItem(dtnav.RootScreenViewers,
		datatugui.MainMenuItem{
			Text:     "Viewers",
			Shortcut: 'v',
			Action:   goViewersScreen,
		})
}
