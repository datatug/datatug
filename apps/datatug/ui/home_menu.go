package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/rivo/tview"
)

type rootScreen int

const (
	loginRootScreen rootScreen = iota
	projectsRootScreen
	settingsRootScreen
)

func newHomeMenu(tui *tapp.TUI, active rootScreen) (menu *homeMenu) {
	newRootScreen := func(newScreen func(tui2 *tapp.TUI) tapp.Screen) func() {
		return func() {
			screen := newScreen(tui)
			tui.SetRootScreen(screen)
		}
	}
	list := menuList().
		AddItem("Login", "", 'l', newRootScreen(newLoginScreen)).
		AddItem("Projects", "", 'p', newRootScreen(newProjectsScreen)).
		AddItem("Settings", "", 's', newRootScreen(newSettingsScreen)).
		AddItem("Exit", "", 'q', func() {
			tui.App.Stop()
		})
	list.SetCurrentItem(int(active))

	defaultBorder(list.Box)

	menu = &homeMenu{
		PanelBase: tapp.NewPanelBase(tui, list, list.Box),
		list:      list,
	}

	return menu
}

var _ tapp.Cell = (*homeMenu)(nil)

type homeMenu struct {
	tapp.PanelBase
	list *tview.List
}
