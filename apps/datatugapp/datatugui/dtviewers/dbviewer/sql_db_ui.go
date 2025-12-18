package dbviewer

import (
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	_ "github.com/mattn/go-sqlite3"
)

func GoSqlDbHome(tui *sneatnav.TUI, dbContext dtviewers.DbContext) error {
	return goTables(tui, sneatnav.FocusToMenu, dbContext)
}
