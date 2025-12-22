package sneatv

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func DefaultBorder(box *tview.Box) {
	box.SetBorder(true)
	box.SetBorderColor(tcell.ColorCornflowerBlue)
	box.SetBorderAttributes(tcell.AttrDim)
	box.SetBorderPadding(1, 0, 1, 1)
}

func SetPanelTitle(box *tview.Box, title string) {
	DefaultBorder(box)
	box.SetTitle(title)
	box.SetTitleAlign(tview.AlignCenter)
	box.SetTitleColor(tview.Styles.TitleColor)
}
