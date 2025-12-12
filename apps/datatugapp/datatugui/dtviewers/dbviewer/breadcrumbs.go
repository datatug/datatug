package dbviewer

import (
	"github.com/datatug/datatug-cli/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
)

func getSqlDbBreadcrumbs(tui *sneatnav.TUI, dbContext dtviewers.DbContext) sneatnav.Breadcrumbs {
	breadcrumbs := dtviewers.GetViewersBreadcrumbs(tui)
	breadcrumbs.Push(sneatv.NewBreadcrumb(dbContext.Driver().ShortTitle, nil))
	return breadcrumbs
}
