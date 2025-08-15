package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
)

type dashboardsSubMenu struct {
	tapp.PanelBase
}

func newDashboardsSidebar(tui *tapp.TUI) *dashboardsSubMenu {
	list := menuList()

	list.AddItem("Add", "", 'a', func() {
		panic("implement me")
	})

	menu := &dashboardsSubMenu{
		PanelBase: tapp.NewPanelBaseFromList(tui, list),
	}

	return menu
}
