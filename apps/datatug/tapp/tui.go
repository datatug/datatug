package tapp

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewTUI(app *tview.Application) *TUI {
	tui := &TUI{
		App:    app,
		Header: NewHeader(),
	}
	app.SetInputCapture(tui.inputCapture)
	return tui
}

func (tui *TUI) inputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch key := event.Key(); key {
	case tcell.KeyRune:
		switch s := string(event.Rune()); s {
		case "q":
			tui.App.Stop()
		default:
			return event
		}
	default:
		return event
	}
	return event
}

type TUI struct {
	App    *tview.Application
	Header *Header
	stack  []Screen
}

func (tui *TUI) StackDepth() int {
	return len(tui.stack)
}

func (tui *TUI) SetRootScreen(screen Screen) {
	tui.stack = []Screen{screen}
	tui.App.SetRoot(screen, screen.Options().FullScreen())
	if err := screen.Activate(); err != nil {
		panic(fmt.Errorf("failed to activate screen: %w", err))
	}
}

func (tui *TUI) PushScreen(screen Screen) {
	tui.stack = append(tui.stack, screen)
	tui.App.SetRoot(screen, screen.Options().FullScreen())
}

func (tui *TUI) PopScreen() {
	for len(tui.stack) > 1 {
		currentScreen := tui.stack[len(tui.stack)-1]
		tui.stack = tui.stack[:len(tui.stack)-1]
		options := currentScreen.Options()
		tui.App.SetRoot(currentScreen, options.fullScreen)
	}
}
