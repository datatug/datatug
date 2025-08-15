package ui

import (
	"github.com/datatug/datatug-cli/pkg/tvprimitives/breadcrumbs"
)

func newHeaderPanel() (bc *breadcrumbs.Breadcrumbs) {
	bc = breadcrumbs.NewBreadcrumbs()
	bc.Push(breadcrumbs.NewBreadcrumb("DataTug", nil))
	bc.Push(breadcrumbs.NewBreadcrumb("Projects", nil))
	bc.Push(breadcrumbs.NewBreadcrumb("Demo project", nil))
	return bc
}
