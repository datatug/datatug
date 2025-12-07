package gcloudui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/rivo/tview"
	"google.golang.org/api/cloudresourcemanager/v3"
)

func newProjectBreadcrumbs(gcProjectCtx CGProjectContext) sneatnav.Breadcrumbs {
	breadcrumbs := newBreadcrumbsProjects(gcProjectCtx.TUI)
	breadcrumbs.Push(sneatv.NewBreadcrumb(gcProjectCtx.Project.DisplayName, func() error {
		return goProject(gcProjectCtx)
	}))
	return breadcrumbs
}

func goProject(gcProjCtx CGProjectContext) error {
	_ = newProjectBreadcrumbs(gcProjCtx)

	menu := newMenuWithProjects(gcProjCtx.TUI, gcProjCtx.Projects)

	list := tview.NewList()
	sneatv.DefaultBorder(list.Box)
	list.SetTitle("Google Cloud Project: " + gcProjCtx.Project.DisplayName)

	list.AddItem("Firestore Database", "", 0, func() {
		_ = goFirestoreDb(gcProjCtx)
	})

	list.AddItem("Firebase Authentication: Users", "", 0, func() {
	})

	content := sneatnav.NewPanelFromList(gcProjCtx.TUI, list)
	gcProjCtx.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToContent))
	return nil
}

func newMenuWithProjects(tui *sneatnav.TUI, projects []*cloudresourcemanager.Project) (menu sneatnav.Panel) {
	list := sneatnav.MainMenuList()
	//sneatv.DefaultBorder(list.Box)
	for _, project := range projects {
		list.AddItem(project.DisplayName, "", 0, func() {})
	}
	return sneatnav.NewPanelFromList(tui, list)
}
