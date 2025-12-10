package sqlviewer

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
)

func GoSqlDB(tui *sneatnav.TUI, filePath string) error {

	return goTables(tui, sneatnav.FocusToMenu, filePath)
}
