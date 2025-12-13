package dtproject

import (
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
)

func goEnvironmentsScreen(ctx ProjectContext, focusTo sneatnav.FocusTo) {

	var menu *projectMenuPanel
	{ // This is too much boilerplate and needs to be simplified
		if existing := ctx.TUI().Menu; existing != nil {
			if m, ok := existing.(*projectMenuPanel); ok {
				menu = m
			}
		}
		if menu == nil {
			menu = newProjectMenuPanel(ctx, "environments")
		}
	}

	//project := ctx.Project()
	//menu.SetProject(project)

	content := newEnvironmentsPanel(ctx)

	ctx.TUI().SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))
}
