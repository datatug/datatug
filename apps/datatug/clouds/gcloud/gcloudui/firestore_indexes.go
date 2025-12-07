package gcloudui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
)

func goFirestoreIndexes(gcProjCtx CGProjectContext) error {
	breadcrumbs := newProjectBreadcrumbs(gcProjCtx)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Firestore", nil))
	menu := firestoreMainMenu(gcProjCtx, firestoreScreenIndexes)

	list := tview.NewList()
	sneatv.DefaultBorder(list.Box)
	list.SetTitle("Firestore Indexes")
	content := sneatnav.NewPanelFromList(gcProjCtx.TUI, list)

	list.AddItem("Loading...", "(not implemented yet)", 0, nil)

	gcProjCtx.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToContent))
	return nil
}
