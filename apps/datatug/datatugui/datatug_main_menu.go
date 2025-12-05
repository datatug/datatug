package datatugui

import (
	"context"

	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
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

func newDataTugMainMenu(tui *sneatnav.TUI, active rootScreen) (menu sneatnav.Panel) {
	handleMenuAction := func(action func(tui2 *sneatnav.TUI, focusTo sneatnav.FocusTo) error) func() {
		return func() {
			if err := action(tui, sneatnav.FocusToContent); err != nil {
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

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		switch index {
		case 0:
			_ = GoHomeScreen(tui, sneatnav.FocusToMenu)
		case 1:
			_ = goProjectsScreen(tui, sneatnav.FocusToMenu)
		case 2:
			_ = goViewersScreen(tui, sneatnav.FocusToMenu)
		case 3:
			_ = goCredentials(tui, sneatnav.FocusToMenu)
		case 4:
			_ = goSettingsScreen(tui, sneatnav.FocusToMenu)
		}
	})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Handle the logic from newDataTugMainMenu: move focus to breadcrumbs when on first item
		switch event.Key() {
		case tcell.KeyUp:
			if list.GetCurrentItem() == 0 {
				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, list)
				return nil
			}
			return event
		case tcell.KeyRight:
			tui.SetFocus(tui.Content)
			return nil
		case tcell.KeyBacktab:
			// Move focus to header (breadcrumbs) when Shift+Tab or Up arrow is pressed on the menu.
			tui.Header.SetFocus(sneatnav.ToBreadcrumbs, list)
			return nil // consume the event
		default:
			return event
		}
	})

	sneatv.DefaultBorder(list.Box)

	return sneatnav.NewPanelFromList(tui, list)
}
