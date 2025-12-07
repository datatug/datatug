package clouds

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
)

type Screen int

const (
	CloudGoogle Screen = iota
	CloudAWS
	CloudAzure
)

func NewCloudsMenu(tui *sneatnav.TUI, active Screen) (menu sneatnav.Panel) {
	list := sneatnav.MainMenuList()
	sneatv.DefaultBorder(list.Box)

	for _, cloud := range registeredClouds {
		list.AddItem(cloud.Name, "", cloud.Shortcut, func() {
			_ = cloud.Action(tui, sneatnav.FocusToMenu)
		})
	}

	list.SetCurrentItem(int(active))

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			list.GetItemSelectedFunc(list.GetCurrentItem())()
		default:
			return event
		}
		return event
	})

	return sneatnav.NewPanelFromList(tui, list)
}
