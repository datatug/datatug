package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/rivo/tview"
)

type projectBasePanel struct {
	project appconfig.ProjectConfig
	tapp.PanelBase
}

func newProjectBasePanel(project appconfig.ProjectConfig, box *tview.Box) projectBasePanel {
	return projectBasePanel{
		project:   project,
		PanelBase: tapp.NewPanelBase(nil, tapp.WithBox(nil, box)),
	}
}
