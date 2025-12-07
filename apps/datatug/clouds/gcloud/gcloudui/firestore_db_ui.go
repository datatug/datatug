package gcloudui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
)

type firestoreScreen int

const (
	firestoreScreenCollections = iota
	firestoreScreenIndexes
)

func goFirestoreDb(gcProjCtx CGProjectContext) error {
	breadcrumbs := newProjectBreadcrumbs(gcProjCtx)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Firestore", nil))
	menu := firestoreMainMenu(gcProjCtx, firestoreScreenCollections)

	list := tview.NewList()
	sneatv.DefaultBorder(list.Box)
	list.SetTitle("Firestore Database")
	content := sneatnav.NewPanelFromList(gcProjCtx.TUI, list)

	gcProjCtx.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToContent))
	return nil
}

func firestoreMainMenu(gcProjCtx CGProjectContext, active firestoreScreen) (menu sneatnav.Panel) {
	list := sneatnav.MainMenuList()
	sneatv.DefaultBorder(list.Box)

	list.AddItem("Collections", "", 0, func() {
		_ = goFirestoreCollections(gcProjCtx)
	})

	list.AddItem("Indexes", "", 0, func() {
		_ = goFirestoreIndexes(gcProjCtx)
	})

	list.SetCurrentItem(int(active))

	return sneatnav.NewPanelFromList(gcProjCtx.TUI, list)
}
