package ui

import (
	"context"
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/strongo/logus"
)

type rootScreen int

const (
	homeRootScreen rootScreen = iota
	projectsRootScreen
	viewersRootScreen
	credentialsRootScreen
	settingsRootScreen
)

func newDataTugMainMenu(tui *tapp.TUI, active rootScreen) (menu *homeMenu) {
	handleMenuAction := func(action func(tui2 *tapp.TUI) error) func() {
		return func() {
			if err := action(tui); err != nil {
				logus.Errorf(context.Background(), "Failed to execute action: %v", err)
				panic(err)
			}
			//tui.SetRootScreen(screen)
		}
	}
	list := menuList().
		AddItem("Home", "", 'h', handleMenuAction(GoHomeScreen)).
		AddItem("Projects", "", 'p', handleMenuAction(goProjectsScreen)).
		AddItem("Viewers", "", 'v', handleMenuAction(goViewersScreen)).
		AddItem("Credentials", "", 'c', handleMenuAction(goCredentials)).
		AddItem("Settings", "", 's', handleMenuAction(goSettingsScreen)).
		AddItem("Exit", "", 'q', func() {
			tui.App.Stop()
		})
	list.SetCurrentItem(int(active))

	list.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		// Handle the logic from newDataTugMainMenu: move focus to breadcrumbs when on first item
		if ev.Key() == tcell.KeyUp {
			if menu.list.GetCurrentItem() == 0 {
				tui.App.SetFocus(tui.Header.Breadcrumbs)
				return nil
			}
		}
		if ev.Key() == tcell.KeyBacktab {
			// Move focus to header (breadcrumbs) when Shift+Tab or Up arrow is pressed on the menu.
			tui.App.SetFocus(tui.Header.Breadcrumbs)
			return nil // consume the event
		}
		return ev
	})

	defaultBorder(list.Box)

	menu = &homeMenu{
		PanelBase: tapp.NewPanelBaseFromList(tui, list),
		list:      list,
	}

	return menu
}

var _ tapp.Cell = (*homeMenu)(nil)

type homeMenu struct {
	tapp.PanelBase
	list *tview.List
}
