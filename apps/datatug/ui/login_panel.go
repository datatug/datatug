package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/rivo/tview"
)

var _ tview.Primitive = (*loginPanel)(nil)
var _ tapp.Cell = (*loginPanel)(nil)

type loginPanel struct {
	tapp.PanelBase
}
