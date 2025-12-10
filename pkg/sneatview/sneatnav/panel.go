package sneatnav

import (
	"github.com/rivo/tview"
)

type Panel interface {
	PrimitiveWithBox
	TakeFocus()
}

type PanelPrimitive interface {
	tview.Primitive
	Box() string
}

//type panelInner interface { // Why we need this?
//	tview.Primitive
//	Box() *tview.Box
//	TakeFocus()
//}

var _ Panel = (*panel[PrimitiveWithBox])(nil)
var _ Cell = (*panel[PrimitiveWithBox])(nil)

type panel[T PrimitiveWithBox] struct {
	PrimitiveWithBox
	tui *TUI
}

func (p panel[T]) TakeFocus() {
	p.tui.App.SetFocus(p.PrimitiveWithBox)
}

func NewPanelWithBoxedPrimitive[T tview.Primitive](tui *TUI, p WithBoxType[T]) Panel {
	return &panel[WithBoxType[T]]{
		PrimitiveWithBox: p,
		tui:              tui,
	}
}

type PanelBase struct {
	tui *TUI
	PrimitiveWithBox
}

func (p PanelBase) TUI() *TUI {
	return p.tui
}

func NewPanelBase(tui *TUI, primitive PrimitiveWithBox) PanelBase {
	if tui == nil {
		panic("tui is nil")
	}
	return PanelBase{tui: tui, PrimitiveWithBox: primitive}
}
