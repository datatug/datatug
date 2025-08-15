package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-cli/pkg/tvprimitives/breadcrumbs"
)

func newProjectsScreen(tui *tapp.TUI) tapp.Screen {
	screen, _ := newDefaultLayout(tui, projectsRootScreen, getProjectsContent)
	tui.Header.Breadcrumbs.Clear()
	tui.Header.Breadcrumbs.Push(breadcrumbs.NewBreadcrumb("Projects", nil))
	return screen
}

var _ tapp.Screen = (*projectsScreen)(nil)

type projectsScreen struct {
	tapp.ScreenBase
	//row *tapp.Row
}
