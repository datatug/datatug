package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
)

type dashboardsSubMenu struct {
	sneatnav.PanelBase
}

func newDashboardsSidebar(tui *sneatnav.TUI) *dashboardsSubMenu {
	list := menuList()

	list.AddItem("Add", "", 'a', func() {
		panic("implement me")
	})

	menu := &dashboardsSubMenu{
		PanelBase: sneatnav.NewPanelBaseFromList(tui, list),
	}

	return menu
}
