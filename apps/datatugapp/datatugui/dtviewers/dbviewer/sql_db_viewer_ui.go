package dbviewer

import (
	"github.com/datatug/datatug/apps/datatugapp/datatugui/dtviewers"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
)

type RecentDB struct {
	Name string
	Path string
}

func goDbViewerSelector(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {

	menu := getDbViewerMenu(tui, focusTo, "SQL DB Viewers")

	breadcrumbs := dtviewers.GetViewersBreadcrumbs(tui)
	breadcrumbs.Push(sneatv.NewBreadcrumb("SQL DB Viewers", nil))

	content := sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(menu, menu.Box))
	tui.SetPanels(tui.Menu, content, sneatnav.WithFocusTo(focusTo))
	return nil
}

func getDbViewerMenu(tui *sneatnav.TUI, focusTo sneatnav.FocusTo, title string) *tview.List {
	list := sneatnav.MainMenuList(tui)
	if title != "" {
		list.SetTitle(title)
	}
	list.AddItem("SQLLite", "", 'l', func() {
		_ = goSqliteHome(tui, focusTo)
	})
	list.AddItem("PostgreSQL", "", 'p', nil)
	setDefaultInputCaptureForList(tui, list)
	return list
}
