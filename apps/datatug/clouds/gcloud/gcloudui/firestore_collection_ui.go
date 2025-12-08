package gcloudui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
)

func goFirestoreCollection(gcProjCtx *CGProjectContext, collectionID string, focusTo sneatnav.FocusTo) error {
	breadcrumbs := firestoreBreadcrumbs(gcProjCtx)
	breadcrumbs.Push(sneatv.NewBreadcrumb(collectionID, nil))

	menu := firestoreMainMenu(gcProjCtx, firestoreScreenCollections, "")

	gcProjCtx.TUI.SetPanels(menu, nil, sneatnav.WithFocusTo(focusTo))

	return nil
}
