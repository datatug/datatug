package sneatnav

import "github.com/rivo/tview"

type PrimitiveWithBox interface {
	tview.Primitive
	GetBox() *tview.Box
}

type withBox struct {
	tview.Primitive
	box *tview.Box
}

func (p withBox) GetBox() *tview.Box {
	return p.box
}

func WithBox(p tview.Primitive, box *tview.Box) PrimitiveWithBox {
	return withBox{
		Primitive: p,
		box:       box,
	}
}
