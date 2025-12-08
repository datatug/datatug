package gcloudui

import (
	"fmt"

	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/cloudresourcemanager/v3"
)

func GoProjects(gcContext *GCloudContext, focusTo sneatnav.FocusTo) error {
	return showProjects(gcContext, focusTo)
}

func OpenProjectsScreen(projects []*cloudresourcemanager.Project) error {
	tui := datatugui.NewDatatugTUI()
	gcContext := &GCloudContext{
		TUI:      tui,
		projects: projects,
	}
	return showProjects(gcContext, sneatnav.FocusToContent)
}

func showProjects(gcContext *GCloudContext, focusTo sneatnav.FocusTo) error {
	breadcrumbs := NewGoogleCloudBreadcrumbs(gcContext)

	breadcrumbs.Push(sneatv.NewBreadcrumb("Projects", func() error {
		return showProjects(gcContext, sneatnav.FocusToContent)
	}))
	menu := newMainMenu(gcContext, ScreenProjects)

	table := tview.NewTable().
		SetSelectable(true, false)
	sneatv.SetPanelTitle(table.Box, "Google Cloud Projects")
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft, tcell.KeyEscape:
			gcContext.TUI.SetFocus(menu)
			return nil
		default:
			return event
		}
	})

	// Header
	headerStyle := tcell.StyleDefault.Bold(true).Reverse(true)

	addHeader := func() {
		table.SetCell(0, 0, tview.NewTableCell("Name").SetSelectable(false).SetStyle(headerStyle))
		table.SetCell(0, 1, tview.NewTableCell("Project ID").SetSelectable(false).SetStyle(headerStyle))
		table.SetCell(0, 2, tview.NewTableCell("#").SetSelectable(false).SetStyle(headerStyle))
	}

	addHeader()
	// Loading row
	table.SetCell(1, 0, tview.NewTableCell("Loading...").SetSelectable(false))

	go func() {
		projects, err := gcContext.GetProjects()
		gcContext.TUI.App.QueueUpdateDraw(func() {
			// Clear rows except header
			table.Clear()
			// Re-add header after Clear
			addHeader()

			if err != nil {
				table.SetCell(1, 0, tview.NewTableCell(fmt.Sprintf("Failed to load projects: %v", err)).SetSelectable(false))
				return
			}
			for i, project := range projects {
				row := i + 1
				gcProjCtx := CGProjectContext{
					GCloudContext: gcContext,
					Project:       project,
				}
				// Store context in the first cell reference
				nameCell := tview.NewTableCell(project.DisplayName).SetReference(gcProjCtx)
				idCell := tview.NewTableCell(project.ProjectId)
				num := ""
				if len(project.Name) > 9 {
					num = project.Name[9:]
				}
				numCell := tview.NewTableCell(num)
				table.SetCell(row, 0, nameCell)
				table.SetCell(row, 1, idCell)
				table.SetCell(row, 2, numCell)
			}
			table.ScrollToBeginning()
		})
	}()

	table.SetSelectedFunc(func(row, column int) {
		if row <= 0 {
			return // header
		}
		cell := table.GetCell(row, 0)
		if cell == nil {
			return
		}
		if ref := cell.GetReference(); ref != nil {
			if ctx, ok := ref.(CGProjectContext); ok {
				_ = goProject(ctx)
			}
		}
	})

	content := sneatnav.NewPanelFromTable(gcContext.TUI, table)

	gcContext.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))

	return nil
}
