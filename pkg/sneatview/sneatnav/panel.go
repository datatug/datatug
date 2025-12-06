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

func NewPanel[T PrimitiveWithBox](tui *TUI, p PrimitiveWithBox) Panel {
	return &panel[T]{
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

func NewPanelBase(tui *TUI,
	primitive PrimitiveWithBox,
	// box *tview.Box, // we have to pass both `list` and `list.Box` as list has no `Box()` required for cell
) PanelBase {
	if tui == nil {
		panic("tui is nil")
	}
	return PanelBase{tui: tui, PrimitiveWithBox: primitive}
}

func NewPanelFromList(tui *TUI, p *tview.List) Panel {
	return NewPanel[withBox[*tview.List]](tui, WithBox(p, p.Box))
}

func NewPanelFromTextView(tui *TUI, p *tview.TextView) Panel {
	return NewPanel[withBox[*tview.TextView]](tui, WithBox(p, p.Box))
}

func NewPanelFromTreeView(tui *TUI, p *tview.TreeView) Panel {
	return NewPanel[withBox[*tview.TreeView]](tui, WithBox(p, p.Box))
}
