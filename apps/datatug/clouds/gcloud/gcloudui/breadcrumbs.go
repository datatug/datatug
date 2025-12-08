package gcloudui

import (
	"github.com/datatug/datatug-cli/apps/datatug/clouds"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
)

func NewGoogleCloudBreadcrumbs(cContext *GCloudContext) sneatnav.Breadcrumbs {
	breadcrumbs := clouds.NewCloudsBreadcrumbs(cContext.TUI)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Google", func() error {
		return GoHome(cContext, sneatnav.FocusToContent)
	}))
	return breadcrumbs
}

func newBreadcrumbsProjects(cContext *GCloudContext) sneatnav.Breadcrumbs {
	breadcrumbs := NewGoogleCloudBreadcrumbs(cContext)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Projects", func() error {
		return GoProjects(cContext, sneatnav.FocusToContent)
	}))
	return breadcrumbs
}
