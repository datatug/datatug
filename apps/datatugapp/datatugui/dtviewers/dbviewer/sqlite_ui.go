package dbviewer

import (
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/rivo/tview"
)

func goSqliteHome(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
	menu := getDbViewerMenu(tui, focusTo, "")
	menuPanel := sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(menu, menu.Box))

	tree := tview.NewTreeView()
	tree.SetTitle("SQLite DB viewer")
	root := tview.NewTreeNode("SQLite DB viewer")
	root.SetSelectable(false)

	tree.SetRoot(root)
	tree.SetTopLevel(1)

	openNode := tview.NewTreeNode("Open SQLite db file")
	root.AddChild(openNode)
	tree.SetCurrentNode(openNode)

	demoNode := tview.NewTreeNode("Demo")
	demoNode.SetSelectable(false)
	root.AddChild(demoNode)

	northwindNode := tview.NewTreeNode(demoDbsFolder + northwindSqliteDbFileName)
	northwindNode.SetSelectedFunc(func() {
		openSqliteDemoDb(northwindSqliteDbFileName)
	})
	demoNode.AddChild(northwindNode)

	content := sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(tree, tree.Box))

	tui.SetPanels(menuPanel, content, sneatnav.WithFocusTo(focusTo))
	return nil
}

const demoDbsFolder = "~/datatug/demo-dbs/"
const northwindSqliteDbFileName = "northwind-sqlite.db"

func openSqliteDemoDb(name string) {
	switch name {
	case northwindSqliteDbFileName:
		_ = downloadFile(
			"https://raw.githubusercontent.com/jpwhite3/northwind-SQLite3/refs/heads/main/dist/northwind.db",
			demoDbsFolder+northwindSqliteDbFileName)

	}
}

func downloadFile(from, to string) error {
	return nil
}
