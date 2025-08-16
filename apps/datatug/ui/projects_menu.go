package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
)

func newProjectsMenu(tui *sneatnav.TUI) *projectsMenu {
	list := menuList()
	list.
		AddItem("Add", "", 'a', nil).
		AddItem("Delete", "", 'd', nil)
	defaultListStyle(list)
	menu := &projectsMenu{
		PanelBase: sneatnav.NewPanelBaseFromList(tui, list),
	}
	setPanelTitle(menu.PanelBase, "")
	return menu
}

type projectsMenu struct {
	sneatnav.PanelBase
}
