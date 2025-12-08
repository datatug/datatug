package gcloudui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newProjectBreadcrumbs(gcProjectCtx CGProjectContext) sneatnav.Breadcrumbs {
	breadcrumbs := newBreadcrumbsProjects(gcProjectCtx.GCloudContext)
	breadcrumbs.Push(sneatv.NewBreadcrumb(gcProjectCtx.Project.DisplayName, func() error {
		return goProject(gcProjectCtx)
	}))
	return breadcrumbs
}

func goProject(gcProjCtx CGProjectContext) error {
	_ = newProjectBreadcrumbs(gcProjCtx)

	menu := newMenuWithProjects(gcProjCtx.GCloudContext)

	list := tview.NewList()
	sneatv.DefaultBorder(list.Box)
	list.SetTitle("Google Cloud Project: " + gcProjCtx.Project.DisplayName)

	list.AddItem("🛢️ Firestore Database", "", 0, func() {
		_ = goFirestoreDb(gcProjCtx)
	})

	list.AddItem("🆔 Firebase Authentication: Users", "", 0, func() {
	})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			gcProjCtx.TUI.SetFocus(menu)
			return nil
		default:
			return event
		}
	})

	content := sneatnav.NewPanelFromList(gcProjCtx.TUI, list)
	gcProjCtx.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToContent))
	return nil
}

func newMenuWithProjects(cContext *GCloudContext) (menu sneatnav.Panel) {
	list := sneatnav.MainMenuList()
	list.SetTitle("Projects")
	sneatv.DefaultBorder(list.Box)
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			cContext.TUI.Header.SetFocus(sneatnav.ToBreadcrumbs, list)
		case tcell.KeyRight:
			cContext.TUI.Content.TakeFocus()
		case tcell.KeyEnter:
			cContext.TUI.Content.TakeFocus()
			cContext.TUI.Content.InputHandler()(event, cContext.TUI.SetFocus)
		default:
			return event
		}
		return event
	})

	projects, err := cContext.GetProjects()
	if err != nil {
		list.AddItem("Failed to load  projects:", err.Error(), 0, nil)
		return sneatnav.NewPanelFromList(cContext.TUI, list)
	}
	for _, project := range projects {
		list.AddItem(project.DisplayName, "", 0, func() {})
	}
	return sneatnav.NewPanelFromList(cContext.TUI, list)
}
