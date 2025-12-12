package dtproject

import (
	"strings"

	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
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
	projectNode.SetSelectedFunc(func() {
		GoProjectScreen(tui, project)
	})

	tree.SetCurrentNode(projectNode)

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		switch key {
		case tcell.KeyUp:
			if tree.GetCurrentNode() == tree.GetRoot() {
				tui.Header.SetFocus(sneatnav.ToBreadcrumbs, tree)
				return nil
			}
			return event
		case tcell.KeyRight, tcell.KeyLeft:
			node := tree.GetCurrentNode()
			if strings.HasPrefix(node.GetText(), "Environments") {
				node.SetExpanded(key == tcell.KeyRight)
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

	dbsNode := tview.NewTreeNode("Databases").SetSelectable(true)
	projectNode.AddChild(dbsNode)
	dbsNode.SetSelectedFunc(func() {
		panic("not implemented")
	})

	envsNode := tview.NewTreeNode("Environments (4)").SetSelectable(true)
	projectNode.AddChild(envsNode)
	envsNode.SetSelectedFunc(func() {
		goEnvironmentsScreen(tui, project)
	})

	envsNode.AddChild(tview.NewTreeNode("Dev"))
	envsNode.AddChild(tview.NewTreeNode("QA"))
	envsNode.AddChild(tview.NewTreeNode("UAT"))
	envsNode.AddChild(tview.NewTreeNode("PROD"))
	envsNode.SetExpanded(false)

	entitiesNode := tview.NewTreeNode("Entities").SetSelectable(true)
	projectNode.AddChild(entitiesNode)

	logsNode := tview.NewTreeNode("Logs").SetSelectable(true)
	projectNode.AddChild(logsNode)

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
