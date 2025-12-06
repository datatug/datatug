package sneatnav

import (
	"github.com/rivo/tview"
)

func MainMenuList() *tview.List {
	list := tview.NewList()
	list.SetWrapAround(false)
	return list
}
