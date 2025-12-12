package clouds

import (
	"github.com/datatug/datatug/pkg/schemers"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
)

type CloudContext struct {
	TUI *sneatnav.TUI
}

type ProjectContext interface {
	Schema() schemers.Provider
}
