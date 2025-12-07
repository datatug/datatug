package gcloudui

import (
	"github.com/datatug/datatug-cli/apps/datatug/clouds"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
)

func NewGoogleCloudBreadcrumbs(tui *sneatnav.TUI) sneatnav.Breadcrumbs {
	breadcrumbs := clouds.NewCloudsBreadcrumbs(tui)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Google", func() error {
		return GoHome(tui, sneatnav.FocusToContent)
	}))
	return breadcrumbs
}

func newBreadcrumbsProjects(tui *sneatnav.TUI) sneatnav.Breadcrumbs {
	breadcrumbs := NewGoogleCloudBreadcrumbs(tui)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Projects", func() error {
		return GoProjects(tui, sneatnav.FocusToContent)
	}))
	return breadcrumbs
}
