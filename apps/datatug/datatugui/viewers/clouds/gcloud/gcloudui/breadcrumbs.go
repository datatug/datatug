package gcloudui

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
)

func NewGoogleCloudBreadcrumbs(cContext *GCloudContext) sneatnav.Breadcrumbs {
	breadcrumbs := viewers.GetViewersBreadcrumbs(cContext.TUI)
	breadcrumbs.Push(sneatv.NewBreadcrumb("Google", func() error {
		return goHome(cContext, sneatnav.FocusToContent)
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
