package gcloudui

import (
	"github.com/datatug/datatug-cli/apps/datatug/clouds"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
)

func NewGoogleCloudBreadcrumbs(gcContext *GCloudContext) sneatnav.Breadcrumbs {
	breadcrumbs := clouds.NewCloudsBreadcrumbs(gcContext.TUI)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Google", func() error {
		return GoHome(gcContext, sneatnav.FocusToContent)
	}))
	return breadcrumbs
}

func newBreadcrumbsProjects(gcContext *GCloudContext) sneatnav.Breadcrumbs {
	breadcrumbs := NewGoogleCloudBreadcrumbs(gcContext)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Projects", func() error {
		return GoProjects(gcContext, sneatnav.FocusToContent)
	}))
	return breadcrumbs
}
