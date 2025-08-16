package sneatnav

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Header struct {
	*tview.Flex
	breadcrumbs *sneatv.Breadcrumbs
	RightMenu   *LoginMenu
}

func (h Header) Breadcrumbs() *sneatv.Breadcrumbs {
	return h.breadcrumbs
}

func NewHeader(tui *TUI, root sneatv.Breadcrumb) *Header {
	header := &Header{
		Flex: tview.NewFlex(),
		breadcrumbs: sneatv.NewBreadcrumbs(root, sneatv.WithInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyDown:
				if tui.focus.from != nil {
					tui.App.SetFocus(tui.focus.from)
					return nil
				}
				return event
			default:
				return event
			}
		})),
		RightMenu: NewLoginMenu(),
	}

	// Calculate fixed width for RightMenu: "(l) Login" = 9 characters + padding
	rightMenuWidth := 11 // "(l) Login" + some padding for borders

	header.AddItem(header.breadcrumbs, 0, 1, true)            // Flexible width (takes remaining space)
	header.AddItem(header.RightMenu, rightMenuWidth, 0, true) // Fixed width

	return header
}
