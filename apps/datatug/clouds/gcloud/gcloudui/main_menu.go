package gcloudui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
)

type Screen int

const (
	ScreenProjects Screen = iota
	ScreenCredentials
)

func newMainMenu(gcContext *GCloudContext, active Screen) (menu sneatnav.Panel) {
	list := sneatnav.MainMenuList()
	sneatv.DefaultBorder(list.Box)

	list.AddItem("Projects", "", 'p', func() {
		_ = GoProjects(gcContext, sneatnav.FocusToContent)
	})
	list.AddItem("Credentials", "", 'c', func() {
		_ = GoCredentials(gcContext, sneatnav.FocusToContent)
	})

	list.SetCurrentItem(int(active))

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			list.GetItemSelectedFunc(list.GetCurrentItem())()
		default:
			return event
		}
		return event
	})

	return sneatnav.NewPanelFromList(gcContext.TUI, list)
}
