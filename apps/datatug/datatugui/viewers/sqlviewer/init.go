package sqlviewer

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
)

const viewerID viewers.ViewerID = "sql"

func RegisterAsViewer() {
	viewers.RegisterViewer(viewers.Viewer{
		ID:       viewerID,
		Name:     "SQL DB viewer",
		Shortcut: '1',
		Action:   goSqlDbHome,
	})
}

func goSqlDbHome(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
	return nil
}
