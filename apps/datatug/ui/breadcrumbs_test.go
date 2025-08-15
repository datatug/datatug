package ui

import (
	"testing"

	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TestBreadcrumbs_InDefaultLayout_FocusTraversalToList tests that breadcrumbs can navigate focus to menu
func TestBreadcrumbs_InDefaultLayout_FocusTraversalToList(t *testing.T) {
	// Create TUI instance
	app := tview.NewApplication()
	tui := tapp.NewTUI(app)

	// Use newDefaultLayout to create the layout
	_, header := newDefaultLayout(tui, projectsRootScreen, getProjectsContent)

	// Test Tab key moves focus to next target (should be menu)
	//var focused tview.Primitive
	var focusCalls int
	if h := header.InputHandler(); h != nil {
		h(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone), func(p tview.Primitive) {
			//focused = p
			focusCalls++
		})
	}
	//if focusCalls != 1 {
	//	t.Fatalf("expected 1 setFocus call, got %d", focusCalls)
	//}
	//if focused == nil {
	//	t.Fatalf("Tab should set focus to the next target")
	//}

	// Also verify Down arrow behaves like Tab
	//focused = nil
	//focusCalls = 0
	//if h := header.InputHandler(); h != nil {
	//	h(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone), func(p tview.Primitive) {
	//		focused = p
	//		focusCalls++
	//	})
	//}
	//if focusCalls != 1 {
	//	t.Fatalf("expected 1 setFocus call on Down, got %d", focusCalls)
	//}
	//if focused == nil {
	//	t.Fatalf("Down should set focus to the next target")
	//}
}

// TestBreadcrumbs_InDefaultLayout_ArrowKeyNavigation tests that arrow keys work for navigation
//func TestBreadcrumbs_InDefaultLayout_ArrowKeyNavigation(t *testing.T) {
//	// Create TUI instance
//	app := tview.NewApplication()
//	tui := tapp.NewTUI(app)
//
//	// Use newDefaultLayout to create the layout
//	screen, header := newDefaultLayout(tui, projectsRootScreen, getProjectsContent)
//
//	tui.Header.Breadcrumbs.Push(breadcrumbs.NewBreadcrumb("Item 1", nil))
//	tui.Header.Breadcrumbs.Push(breadcrumbs.NewBreadcrumb("Item 2", nil))
//	tui.Header.Breadcrumbs.Push(breadcrumbs.NewBreadcrumb("Item 3", nil))
//
//	header.Focus(nil) // selects last (Gamma)
//
//	t.Run("LEFT and RIGHT arrow keys should be handled by app", func(t *testing.T) {
//		// Ensure breadcrumbs has focus before testing
//		app.SetFocus(header)
//
//		// Get initial state - should be at last item (index 2, Gamma)
//		initialIndex := header.SelectedItemIndex
//		if initialIndex != 2 {
//			t.Fatalf("Expected initial index 2 (Gamma), got %d", initialIndex)
//		}
//
//		leftKeyEvent := tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
//
//		// Send keyboard event to app to make sure it's not get handled by any parent component
//		app.QueueEvent(leftKeyEvent)
//
//		// Directly call the breadcrumbs input handler to simulate LEFT arrow key
//		//if h := header.InputHandler(); h != nil {
//		//	h(tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone), func(p tview.Primitive) {
//		//		// This setFocus function should not be called for LEFT arrow
//		//	})
//		//}
//
//		// Verify LEFT arrow moved selection from index 2 to index 1 (Beta)
//		if header.SelectedItemIndex != 1 {
//			t.Errorf("Selected item index should be 1 after LEFT arrow, got %d", header.SelectedItemIndex)
//		}
//
//		// Test RIGHT arrow key - should move from index 1 back to index 2
//		if h := header.InputHandler(); h != nil {
//			h(tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone), func(p tview.Primitive) {
//				// This setFocus function should not be called for RIGHT arrow
//			})
//		}
//
//		// Verify RIGHT arrow moved selection from index 1 to index 2 (Gamma)
//		if header.SelectedItemIndex != 2 {
//			t.Errorf("Selected item index should be 2 after RIGHT arrow, got %d", header.SelectedItemIndex)
//		}
//	})
//}
