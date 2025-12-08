package gcloudui

import (
	"github.com/datatug/datatug-cli/apps/datatug/clouds"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
)

type firestoreScreen int

const (
	firestoreScreenCollections = iota
	firestoreScreenIndexes
)

func goFirestoreDb(gcProjCtx CGProjectContext) error {
	breadcrumbs := newProjectBreadcrumbs(gcProjCtx)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Firestore", nil))
	menu := clouds.NewCloudsMenu(gcProjCtx.TUI, clouds.CloudGoogle)

	content := firestoreMainMenu(gcProjCtx, firestoreScreenCollections, "Firestore Database")

	gcProjCtx.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToContent))
	return nil
}

func firestoreMainMenu(gcProjCtx CGProjectContext, active firestoreScreen, title string) (menu sneatnav.Panel) {
	list := sneatnav.MainMenuList()
	sneatv.DefaultBorder(list.Box)

	list.AddItem("Collections", "", 0, func() {
		_ = goFirestoreCollections(gcProjCtx)
	})

	list.AddItem("Indexes", "", 0, func() {
		_ = goFirestoreIndexes(gcProjCtx)
	})

	list.SetCurrentItem(int(active))

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRight {
			gcProjCtx.TUI.SetFocus(gcProjCtx.TUI.Content)
		}
		return event
	})

	return sneatnav.NewPanelFromList(gcProjCtx.TUI, list)
}
