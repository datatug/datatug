package dtproject

import (
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/rivo/tview"
)

func newEnvironmentsPanel(tui *sneatnav.TUI, _ *appconfig.ProjectConfig) sneatnav.Panel {
	list := tview.NewList()
	list.SetWrapAround(false)
	list.AddItem("DEV", "Development", 'd', nil)
	list.AddItem("QA", "Quality Assurance", 'q', nil)
	list.AddItem("UAT", "User Acceptance Testing", 'u', nil)
	list.AddItem("PROD", "Production", 'p', nil)

	return sneatnav.NewPanelWithBoxedPrimitive(tui, sneatnav.WithBox(list, list.Box))
}
