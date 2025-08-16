package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/rivo/tview"
)

type projectBasePanel struct {
	project appconfig.ProjectConfig
	sneatnav.PanelBase
}

func newProjectBasePanel(project appconfig.ProjectConfig, box *tview.Box) projectBasePanel {
	return projectBasePanel{
		project:   project,
		PanelBase: sneatnav.NewPanelBase(nil, sneatnav.WithBox(nil, box)),
	}
}
