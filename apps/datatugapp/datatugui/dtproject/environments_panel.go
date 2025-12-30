package dtproject

import (
	"github.com/datatug/datatug-core/pkg/datatug"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
)

func newEnvironmentsPanel(ctx ProjectContext) sneatnav.Panel {
	project := ctx.Project()
	environments, err := project.GetEnvironments(ctx)
	return newListPanel(ctx.TUI(), "Environments", environments, func(e *datatug.Environment) (string, string) {
		return e.ID, e.Title
	}, err)
}
