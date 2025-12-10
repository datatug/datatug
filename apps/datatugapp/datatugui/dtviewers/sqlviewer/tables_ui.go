package sqlviewer

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func goTables(tui *sneatnav.TUI, _ sneatnav.FocusTo, filePath string) error {

	breadcrumbs := getSqlDbBreadcrumbs(tui, filePath)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Tables", nil))

	menu := newSqlDbMenu(tui, SqlDbScreenTables, filePath)

	textView := tview.NewTextView()
	textView.SetTitle("Tables @ " + filePath)
	setDefaultInputCapture(tui, textView)

	content := sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(textView, textView.Box))

	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToContent))

	return nil
}

func goViews(tui *sneatnav.TUI, _ sneatnav.FocusTo, filePath string) error {
	breadcrumbs := getSqlDbBreadcrumbs(tui, filePath)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Views", nil))

	menu := newSqlDbMenu(tui, SqlDbScreenViews, filePath)

	textView := tview.NewTextView()
	textView.SetTitle("Views @ " + filePath)
	setDefaultInputCapture(tui, textView)

	content := sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(textView, textView.Box))

	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToContent))
	return nil
}

func setDefaultInputCapture(tui *sneatnav.TUI, c interface {
	tview.Primitive
	SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *tview.Box
}) {
	c.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			tui.Menu.TakeFocus()
			return nil
		case tcell.KeyUp:
			tui.Header.SetFocus(sneatnav.ToBreadcrumbs, c)
			return nil
		default:
			return event
		}
	})
}
