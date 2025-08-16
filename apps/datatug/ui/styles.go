package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func defaultBorder(box *tview.Box) {
	box.SetBorder(true)
	box.SetBorderColor(tcell.ColorCornflowerBlue)
	box.SetBorderAttributes(tcell.AttrDim)
	box.SetBorderPadding(1, 0, 1, 1)
}

func defaultListStyle(list *tview.List) {
	list.SetWrapAround(false)
	//defaultBorder(list.Box)
}

func setPanelTitle(panel sneatnav.PanelBase, title string) {
	box := panel.GetBox()
	defaultBorder(box)
	box.SetTitle(title)
	box.SetTitleAlign(tview.AlignCenter)
	box.SetTitleColor(tview.Styles.TitleColor)
}
