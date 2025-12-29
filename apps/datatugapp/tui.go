package datatugapp

import (
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtproject"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
)

func NewDatatugTUI() (tui *sneatnav.TUI) {
	app := tview.NewApplication()
	app.EnableMouse(true)

	tui = sneatnav.NewTUI(app, sneatv.NewBreadcrumb(" ⛴ DataTug", func() error {
		return goProjectScreen(tui, sneatnav.FocusToMenu)
	}))
	return
}

var goProjectScreen = dtproject.GoProjectsScreen
