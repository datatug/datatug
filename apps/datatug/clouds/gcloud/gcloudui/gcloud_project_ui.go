package gcloudui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
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

	list.AddItem("Firestore Database", "", 0, func() {
		_ = goFirestoreDb(gcProjCtx)
	})

	list.AddItem("Firebase Authentication: Users", "", 0, func() {
	})

	content := sneatnav.NewPanelFromList(gcProjCtx.TUI, list)
	gcProjCtx.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToContent))
	return nil
}

func newMenuWithProjects(gcContext *GCloudContext) (menu sneatnav.Panel) {
	list := sneatnav.MainMenuList()
	//sneatv.DefaultBorder(list.Box)
	projects, err := gcContext.GetProjects()
	if err != nil {
		list.AddItem("Failed to load  projects:", err.Error(), 0, nil)
		return sneatnav.NewPanelFromList(gcContext.TUI, list)
	}
	for _, project := range projects {
		list.AddItem(project.DisplayName, "", 0, func() {})
	}
	return sneatnav.NewPanelFromList(gcContext.TUI, list)
}
