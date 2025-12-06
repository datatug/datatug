package fsviewer

import (
	"github.com/datatug/datatug-cli/apps/datatug/dtscreeens/dtviewers"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
)

//var viewersList = tview.NewList()
//
//var AddViewer = viewersList.AddItem

func RegisterModule(tui *sneatnav.TUI) {

	// Register as a viewer to the list of DataTug viewers
	dtviewers.AddViewer("Firestore viewer", "Browse & edit data in Firestore databases", '1', func() {
		_ = OpenFirestoreViewer(tui, sneatnav.FocusToMenu)
	})
}
