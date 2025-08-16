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

var _ Panel = (*PanelBase)(nil)
var _ Cell = (*PanelBase)(nil)

type PanelBase struct {
	tui *TUI
	PrimitiveWithBox
}

func (p *PanelBase) TakeFocus() {
	//return p.box.
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

func NewPanelBaseFromList(tui *TUI, p *tview.List) PanelBase {
	return NewPanelBase(tui, WithBox(p, p.Box))
}

func NewPanelBaseFromTextView(tui *TUI, p *tview.TextView) PanelBase {
	return NewPanelBase(tui, WithBox(p, p.Box))
}
