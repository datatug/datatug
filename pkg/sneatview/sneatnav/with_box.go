package sneatnav

import "github.com/rivo/tview"

type PrimitiveWithBox interface {
	tview.Primitive
	GetBox() *tview.Box
}

var _ PrimitiveWithBox = (*withBox[tview.Primitive])(nil)

type withBox[T tview.Primitive] struct {
	tview.Primitive
	box *tview.Box
}

func (p withBox[T]) GetBox() *tview.Box {
	return p.box
}
func (p withBox[T]) GetPrimitive() T {
	return p.Primitive.(T)
}

func WithBox[T tview.Primitive](p T, box *tview.Box) PrimitiveWithBox {
	return withBox[T]{
		Primitive: p,
		box:       box,
	}
}
