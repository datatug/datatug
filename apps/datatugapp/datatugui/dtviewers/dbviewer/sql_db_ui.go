package dbviewer

import (
	"errors"

	"github.com/datatug/datatug-cli/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	_ "github.com/mattn/go-sqlite3"
)

func GoDbViewerHome(tui *sneatnav.TUI, dbContext dtviewers.DbContext) error {
	if dbContext == nil {
		return errors.New("not implemented yet - dbContext is nil")
	}
	return goTables(tui, sneatnav.FocusToMenu, dbContext)
}
