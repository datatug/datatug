package gcloudui

import (
	"context"
	"fmt"
	"slices"

	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/iterator"
)

func goFirestoreCollection(gcProjCtx *CGProjectContext, collectionID string, focusTo sneatnav.FocusTo) error {
	breadcrumbs := firestoreBreadcrumbs(gcProjCtx)
	breadcrumbs.Push(sneatv.NewBreadcrumb(collectionID, nil))

	menu := firestoreMainMenu(gcProjCtx, firestoreScreenCollections, "")

	// Content: table with first 100 docs
	// Enable cell selection (row & column) so only a single cell is highlighted
	table := tview.NewTable().SetFixed(1, 0).SetSelectable(true, true)
	sneatv.DefaultBorder(table.Box)

	title := "Collection: " + collectionID
	if gcProjCtx.Project != nil && gcProjCtx.Project.ProjectId != "" {
		title += " — " + gcProjCtx.Project.ProjectId
	}
	table.SetTitle(title)

	// Header
	headerStyle := tview.Styles.SecondaryTextColor
	table.SetCell(0, 0, tview.NewTableCell("Doc ID").SetTextColor(headerStyle).SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("Data (JSON)").SetTextColor(headerStyle).SetSelectable(false))

	// Loading placeholder
	table.SetCell(1, 0, tview.NewTableCell("Loading...").SetSelectable(false))

	content := sneatnav.NewPanelFromTable(gcProjCtx.TUI, table)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft {
			// Only switch focus to the menu if we're in the first column (# / Doc ID)
			_, col := table.GetSelection()
			if col == 0 {
				gcProjCtx.TUI.SetFocus(menu)
				return nil
			}
			// Otherwise, let the table handle the left navigation within cells
		}
		return event
	})

	// Async load first 100 docs
	go func() {
		ctx := context.Background()
		projectID := ""
		if gcProjCtx.Project != nil {
			projectID = gcProjCtx.Project.ProjectId
		}
		client, err := newFirestoreClient(ctx, projectID)
		if err != nil {
			gcProjCtx.TUI.App.QueueUpdateDraw(func() {
				// Clear and show error
				table.Clear()
				table.SetCell(0, 0, tview.NewTableCell("Error").SetSelectable(false))
				table.SetCell(1, 0, tview.NewTableCell(err.Error()).SetSelectable(false))
				table.SetCell(2, 0, tview.NewTableCell("Hint: re-login or check scopes in Firestore Collections screen").SetSelectable(false))
			})
			return
		}
		defer func() { _ = client.Close() }()

		iter := client.Collection(collectionID).Limit(100).Documents(ctx)
		type row struct {
			id   string
			data map[string]any
		}
		var rows []row
		var columns []string
		for {
			snap, err := iter.Next()
			if err != nil {
				if err == iterator.Done {
					break
				}
				gcProjCtx.TUI.App.QueueUpdateDraw(func() {
					table.Clear()
					table.SetCell(0, 0, tview.NewTableCell("Error").SetSelectable(false))
					table.SetCell(1, 0, tview.NewTableCell(err.Error()).SetSelectable(false))
				})
				return
			}
			r := row{
				id:   snap.Ref.ID,
				data: snap.Data(),
			}
			for col := range r.data {
				if !slices.Contains(columns, col) {
					columns = append(columns, col)
				}
			}
			//b, _ := json.Marshal(snap.Data())
			rows = append(rows, r)
		}

		slices.Sort(columns)

		gcProjCtx.TUI.App.QueueUpdateDraw(func() {
			const maxWidth = 15
			table.Clear()
			// Recreate header
			table.SetCell(0, 0, tview.NewTableCell("#").SetTextColor(headerStyle).SetSelectable(false))
			for i, col := range columns {
				cell := tview.NewTableCell(col).SetTextColor(headerStyle).SetSelectable(false)
				cell.MaxWidth = maxWidth
				table.SetCell(0, i+1, cell)
			}
			//table.SetCell(0, 1, tview.NewTableCell("Data (JSON)").SetTextColor(headerStyle).SetSelectable(false))
			if len(rows) == 0 {
				table.SetCell(1, 0, tview.NewTableCell("No documents").SetSelectable(false))
				return
			}
			for i, r := range rows {
				table.SetCell(i+1, 0, tview.NewTableCell(r.id))
				for j, col := range columns {
					if v, hasVal := r.data[col]; hasVal {
						cell := tview.NewTableCell(fmt.Sprintf("%v", v))
						cell.MaxWidth = maxWidth
						table.SetCell(i+1, j+1, cell)
					}
				}
			}
		})
	}()

	gcProjCtx.TUI.SetPanels(menu, content, sneatnav.WithFocusTo(focusTo))

	return nil
}
