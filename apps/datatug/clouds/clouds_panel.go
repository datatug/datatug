package clouds

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
	"github.com/datatug/datatug-cli/apps/datatug/dtnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func goClouds(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
	breadcrumbs := tui.Header.Breadcrumbs()
	breadcrumbs.Clear()
	breadcrumbs.Push(sneatv.NewBreadcrumb("Clouds", nil))
	menu := datatugui.NewDataTugMainMenu(tui, dtnav.RootScreenClouds)
	content := newCloudsPanel(tui)
	tui.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))
	return nil
}

func newCloudsPanel(tui *sneatnav.TUI) *cloudsPanel {
	//tree := tview.NewTreeView()

	list := tview.NewList()

	for _, cloud := range registeredClouds {
		list.AddItem(cloud.Name, "", cloud.Shortcut, func() {
			if err := cloud.Action(tui, sneatnav.FocusToContent); err != nil {
				panic(err) // TODO: Show error to user
			}
		})
	}

	//cloudsNode := tview.NewTreeNode("Big Clouds").SetSelectable(false)
	//tree.SetRoot(cloudsNode)
	//
	//cloudsNode.AddChild(tview.NewTreeNode("Google Cloud").SetSelectable(true))
	//cloudsNode.AddChild(tview.NewTreeNode("Amazon Web Services").SetSelectable(true))
	//cloudsNode.AddChild(tview.NewTreeNode("Microsoft Azure").SetSelectable(true))
	//
	//tree.SetCurrentNode(cloudsNode.GetChildren()[0])

	panel := &cloudsPanel{
		PanelBase: sneatnav.NewPanelBase(tui, sneatnav.WithBox(list, list.Box)),
		//tree:      tree,
		list: list,
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC, tcell.KeyLeft:
			tui.SetFocus(tui.Menu)
			return nil
		default:
			return event
		}
	})
	sneatv.SetPanelTitle(panel.GetBox(), "Clouds")
	return panel
}

type cloudsPanel struct {
	sneatnav.PanelBase
	//tree *tview.TreeView
	list *tview.List
}

func (p *cloudsPanel) TakeFocus() {
	p.TUI().SetFocus(p.list)
}
