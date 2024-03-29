package ui

import (
	"github.com/datatug/datatug/packages/cli/config"
	"github.com/rivo/tview"
)

type dashboardsPanel struct {
	projectBasePanel
}

func newDashboardsPanel(project config.ProjectConfig) *dashboardsPanel {

	content := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("List of dashboards here")

	defaultBorder(content.Box)

	return &dashboardsPanel{
		projectBasePanel: newProjectBasePanel(project, content, content.Box),
	}
}
