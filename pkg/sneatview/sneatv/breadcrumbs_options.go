package sneatv

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func WithSeparator(separator string) func(bc *Breadcrumbs) {
	return func(bc *Breadcrumbs) {
		bc.separator = separator
	}
}

type InputHandler = func(event *tcell.EventKey, setFocus func(p tview.Primitive)) *tcell.EventKey

func WithInputHandler(inputHandler InputHandler) func(bc *Breadcrumbs) {
	return func(bc *Breadcrumbs) {
		bc.inputHandler = inputHandler
	}
}
