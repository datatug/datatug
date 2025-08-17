package datatugui

import (
	"github.com/rivo/tview"
)

func menuList() *tview.List {
	list := tview.NewList()
	list.SetWrapAround(false)
	return list
}
