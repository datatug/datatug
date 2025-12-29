package sneatnav

import (
	"errors"
	"fmt"

	"github.com/rivo/tview"
)

type ActionMenuItem struct {
	ID           string
	Title        string
	SelectedFunc func()
}

func newActionsMenu(app *tview.Application) ActionsMenu {
	am := ActionsMenu{app: app, flex: tview.NewFlex()}
	am.Clear()
	return am
}

type ActionsMenu struct {
	app       *tview.Application
	flex      *tview.Flex
	menuItems []ActionMenuItem
}

func (m *ActionsMenu) Clear() {
	m.menuItems = nil
	_ = m.RegisterActionMenuItems(
		ActionMenuItem{
			ID: "Quit", Title: "(Q)uit",
			SelectedFunc: func() {
				m.app.Stop()
			},
		},
		ActionMenuItem{
			ID: "Help", Title: "F1 - Help",
			SelectedFunc: func() {
				return
			},
		},
	)
}

func (m *ActionsMenu) RegisterActionMenuItems(items ...ActionMenuItem) error {
	for i, item := range items {
		if item.ID == "" {
			panic(fmt.Sprintf("items[%d] has no ID, title=%s", i, item.Title))
		}
		for _, menuItem := range m.menuItems {
			if menuItem.ID == item.ID {
				return errors.New("and attempt to register a menu item with already registered ID=" + item.ID)
			}
		}
		m.menuItems = append(m.menuItems, item)
		title := item.Title
		if title == "" {
			title = item.ID
		}
		b := tview.NewButton(item.Title)
		b.SetSelectedFunc(item.SelectedFunc)
		m.flex.AddItem(b, len(item.Title)+2, 0, false)
		m.flex.AddItem(nil, 1, 0, false) // right margin
	}
	return nil
}
