package dtproject

import (
	"github.com/datatug/datatug-core/pkg/datatug"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
)

func goDatabasesScreen(ctx ProjectContext, focusTo sneatnav.FocusTo) {

	menu := getOrCreateProjectMenuPanel(ctx, "environments")

	//project := ctx.Project()
	//menu.SetProject(project)

	content := newDatabasesPanel(ctx)

	ctx.TUI().SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))
}

func newDatabasesPanel(ctx ProjectContext) sneatnav.Panel {
	project := ctx.Project()
	dbServers, err := project.GetDbServers(ctx)
	return newListPanel(ctx.TUI(), "Databases", dbServers, func(s *datatug.ProjDbServer) (string, string) {
		return s.ID, s.Title
	}, err)
}
