package sqlviewer

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/dtviewers"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
)

const viewerID dtviewers.ViewerID = "sql"

func RegisterAsViewer() {
	dtviewers.RegisterViewer(dtviewers.Viewer{
		ID:       viewerID,
		Name:     "SQL DB viewer",
		Shortcut: '1',
		Action:   goSqlDbHome,
	})
}

func goSqlDbHome(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
	breadcrumbs := dtviewers.GetViewersBreadcrumbs(tui)
	breadcrumbs.Push(sneatv.NewBreadcrumb("SQL", nil))

	textView := tview.NewTextView()
	sneatv.DefaultBorder(textView.Box)
	textView.SetTitle("SQL DB Viewer")
	content := sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(textView, textView.Box))
	tui.SetPanels(nil, content, sneatnav.WithFocusTo(focusTo))

	return nil
}
