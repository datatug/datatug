package clouds

import (
	"github.com/datatug/datatug-cli/pkg/dbschema"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
)

type CloudContext struct {
	TUI *sneatnav.TUI
}

type ProjectContext interface {
	Schema() dbschema.Provider
}
