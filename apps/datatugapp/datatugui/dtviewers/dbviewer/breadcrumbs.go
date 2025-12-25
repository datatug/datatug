package dbviewer

import (
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
)

func getSqlDbBreadcrumbs(tui *sneatnav.TUI, dbContext dtviewers.DbContext) sneatnav.Breadcrumbs {
	breadcrumbs := GetDbViewersBreadcrumbs(tui)
	driverBreadcrumb := sneatv.NewBreadcrumb(dbContext.Driver().ShortTitle, func() error {
		return goSqliteHome(tui, sneatnav.FocusToContent)
	})
	breadcrumbs.Push(driverBreadcrumb)

	if name := dbContext.Name(); name != "" {
		dbBreadcrumb := sneatv.NewBreadcrumb(name, func() error {
			return GoSqlDbHome(tui, dbContext)
		})
		breadcrumbs.Push(dbBreadcrumb)
	}
	return breadcrumbs
}
