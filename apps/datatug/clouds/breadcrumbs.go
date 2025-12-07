package clouds

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
)

func NewCloudsBreadcrumbs(tui *sneatnav.TUI) sneatnav.Breadcrumbs {
	breadcrumbs := tui.Header.Breadcrumbs()
	breadcrumbs.Clear()
	breadcrumbs.Push(sneatv.NewBreadcrumb("Clouds", func() error {
		return goClouds(tui, sneatnav.FocusToContent)
	}))
	return breadcrumbs
}
