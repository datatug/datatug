package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
)

func newProjectsMenu(tui *tapp.TUI) *projectsMenu {
	list := menuList()
	list.
		AddItem("Add", "", 'a', nil).
		AddItem("Delete", "", 'd', nil)
	defaultListStyle(list)
	menu := &projectsMenu{
		PanelBase: tapp.NewPanelBase(tui, list, list.Box),
	}
	return menu
}

type projectsMenu struct {
	tapp.PanelBase
}
