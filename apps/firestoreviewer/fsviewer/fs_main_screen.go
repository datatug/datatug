package fsviewer

import (
	"github.com/datatug/datatug-cli/apps/datatug/dtscreeens/dtviewers"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
)

func OpenFirestoreViewer(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {

	breadcrumbs := dtviewers.GetViewersBreadcrumbs(tui)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Firebase", nil))
	breadcrumbs.Push(sneatv.NewBreadcrumb("Firestore", nil))
	list := tview.NewList()

	// Add the two required items
	list.AddItem("Firebase projects", "Browse & edit data in Firestore databases", '1', nil)

	menu := newFirestoreViewerMainMenu(tui, 0)
	content := sneatnav.NewPanelFromList(tui, list)

	tui.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))

	return nil
}
