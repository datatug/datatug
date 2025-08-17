package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
)

func newProjectsMenuPanel(tui *sneatnav.TUI) sneatnav.Panel {
	list := menuList()
	list.
		AddItem("Add", "", 'a', nil).
		AddItem("Delete", "", 'd', nil)
	defaultListStyle(list)
	setPanelTitle(list.Box, "")
	return sneatnav.NewPanelFromList(tui, list)
}
