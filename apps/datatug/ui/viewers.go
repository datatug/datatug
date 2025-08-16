package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-cli/pkg/tvprimitives/breadcrumbs"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var _ tview.Primitive = (*viewersPanel)(nil)
var _ tapp.Cell = (*viewersPanel)(nil)

type viewersPanel struct {
	tapp.PanelBase
	list *tview.List
}

func (p *viewersPanel) Draw(screen tcell.Screen) {
	p.list.Draw(screen)
}

func goViewersScreen(tui *tapp.TUI) error {
	tui.Header.Breadcrumbs.Clear()
	tui.Header.Breadcrumbs.Push(breadcrumbs.NewBreadcrumb("Viewers", nil))
	list := tview.NewList()

	// Add the two required items
	list.AddItem("Firestore viewer", "Browse & edit data in Firestore databases", '1', nil)
	list.AddItem("SQL DB viewer", "Browse & query SQL databases", '2', nil)

	// Set secondary text color to gray
	list.SetSecondaryTextColor(tcell.ColorDarkGray)

	content := &viewersPanel{
		PanelBase: tapp.NewPanelBaseFromList(tui, list),
		list:      list,
	}

	defaultBorder(content.list.Box)
	// Set spacing between items to 1 line by increasing vertical padding
	content.list.SetBorderPadding(1, 1, 1, 1)
	content.list.SetTitle(" Viewers ")
	content.list.SetTitleAlign(tview.AlignLeft)

	menu := newDataTugMainMenu(tui, viewersRootScreen)
	tui.SetPanels(menu, content)
	return nil
}
