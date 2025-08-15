package tapp

import (
	"github.com/datatug/datatug-cli/pkg/tvprimitives/breadcrumbs"
	"github.com/rivo/tview"
)

type Header struct {
	*tview.Flex
	Breadcrumbs *breadcrumbs.Breadcrumbs
	RightMenu   *tview.TextView
}

func NewHeader() *Header {
	header := &Header{
		Flex:        tview.NewFlex(),
		Breadcrumbs: breadcrumbs.NewBreadcrumbs(breadcrumbs.NewBreadcrumb(" ⛴ DataTug", nil)),
		RightMenu:   tview.NewTextView().SetText("Sign In"),
	}

	header.AddItem(header.Breadcrumbs, 0, 1, true)
	header.AddItem(header.RightMenu, 8, 1, true)

	return header
}
