package sqlviewer

import (
	"github.com/datatug/datatug-cli/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug-cli/pkg/dtio"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
)

func getSqlDbBreadcrumbs(tui *sneatnav.TUI, filePath string) sneatnav.Breadcrumbs {
	breadcrumbs := dtviewers.GetViewersBreadcrumbs(tui)
	if dtio.IsSQLite(filePath) {
		breadcrumbs.Push(sneatv.NewBreadcrumb("SQLite", nil))
	} else {
		breadcrumbs.Push(sneatv.NewBreadcrumb("SQL DB", nil))
	}
	return breadcrumbs
}
