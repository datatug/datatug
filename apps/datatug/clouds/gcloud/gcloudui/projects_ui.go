package gcloudui

import (
	"context"
	"fmt"

	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
	"github.com/datatug/datatug-cli/pkg/auth/gauth"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/cloudresourcemanager/v3"
)

func GoProjects(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
	ctx := context.Background()
	projects, err := gauth.GetGCloudProjects(ctx)
	if err != nil {
		return err
	}
	return showProjects(GCloudContext{
		TUI:      tui,
		Projects: projects,
	}, focusTo)
}

func OpenProjectsScreen(projects []*cloudresourcemanager.Project) error {
	tui := datatugui.NewDatatugTUI()
	gcContext := GCloudContext{
		TUI:      tui,
		Projects: projects,
	}
	return showProjects(gcContext, sneatnav.FocusToContent)
}

func showProjects(gcContext GCloudContext, focusTo sneatnav.FocusTo) error {
	breadcrumbs := NewGoogleCloudBreadcrumbs(gcContext.TUI)

	breadcrumbs.Push(sneatv.NewBreadcrumb("Projects", func() error {
		return showProjects(gcContext, sneatnav.FocusToContent)
	}))
	menu := newMainMenu(gcContext.TUI, ScreenProjects)

	list := tview.NewList()
	sneatv.SetPanelTitle(list.Box, "Google Cloud Projects")
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft, tcell.KeyEscape:
			gcContext.TUI.SetFocus(menu)
			return nil
		default:
			return event
		}
	})

	for _, project := range gcContext.Projects {
		gcProjCtx := CGProjectContext{
			GCloudContext: gcContext,
			Project:       project,
		}
		list.AddItem(
			project.DisplayName,
			fmt.Sprintf("%s (# %s)", project.ProjectId, project.Name[9:]),
			0,
			func() {
				_ = goProject(gcProjCtx)
			},
		)
	}

	content := sneatnav.NewPanelFromList(gcContext.TUI, list)

	gcContext.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))

	return nil
}
