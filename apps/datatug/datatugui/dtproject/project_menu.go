package dtproject

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ProjectScreenID string

const (
	ProjectScreenDashboards   = "dashboards"
	ProjectScreenEnvironments = "environments"
)

func newProjectMenuPanel(tui *sneatnav.TUI, project *appconfig.ProjectConfig, currentScreen ProjectScreenID) sneatnav.Panel {
	tree := tview.NewTreeView()

	//datatugNode := tview.NewTreeNode("DataTug").SetSelectable(true)
	//tree.SetRoot(datatugNode)
	//
	//projectsNode := tview.NewTreeNode("Projects").SetSelectable(true)
	//datatugNode.AddChild(projectsNode)

	//projectTitle := " 📁 " + GetProjectShortTitle(project) // TODO: emoji breaks borders
	projectTitle := GetProjectShortTitle(project)
	projectNode := tview.NewTreeNode(projectTitle).SetSelectable(true)
	tree.SetRoot(projectNode)

	tree.SetCurrentNode(projectNode)

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			if tree.GetCurrentNode() == tree.GetRoot() {
				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, tree)
				return nil
			}
			return event
		default:
			return event
		}
	})

	dashboardsNode := tview.NewTreeNode("Dashboards").SetSelectable(true)
	projectNode.AddChild(dashboardsNode)
	dashboardsNode.SetSelectedFunc(func() {
		goProjectDashboards(tui, project)
	})
	if currentScreen == ProjectScreenDashboards {
		tree.SetCurrentNode(dashboardsNode)
	}

	envsNode := tview.NewTreeNode("Environments").SetSelectable(true)
	projectNode.AddChild(envsNode)
	envsNode.SetSelectedFunc(func() {
		goEnvironmentsScreen(tui, project)
	})

	if currentScreen == ProjectScreenEnvironments {
		tree.SetCurrentNode(envsNode)
	}

	/*
		list := tview.NewList().
			//AddItem("Databases", "", 'D', nil).
			AddItem("Dashboards", "", 'B', func() {
				goProjectDashboards(tui, project)
			}).
			AddItem("Environments", "", 'E', func() {
				goEnvironmentsScreen(tui, project)
			})

		AddItem("Queries", "", 'Q', nil).
		AddItem("Web UI", "", 'W', nil)

		currentItem := -1
		switch currentScreen {
		case ProjectScreenDashboards:
			currentItem = 0
		case ProjectScreenEnvironments:
			currentItem = 1
		}
		if currentItem >= 0 {
			list.SetCurrentItem(currentItem)
		}

		sneatv.DefaultListStyle(list)
	*/

	sneatv.DefaultBorder(tree.Box)

	return sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(tree, tree.Box))
}
