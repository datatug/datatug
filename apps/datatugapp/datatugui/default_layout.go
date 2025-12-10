package datatugui

//func newDefaultLayout(
//	tui *sneatnav.TUI, selectedMenuItem RootScreen, getContent func(tui *sneatnav.TUI) (sneatnav.Panel, error),
//) sneatnav.Screen {
//	addMainRow(tui, selectedMenuItem, tui.Grid, getContent)
//
//	return nil
//}

//func addMainRow(
//	tui *sneatnav.TUI, selectedMenuItem RootScreen, grid *tview.Grid,
//	getContent func(tui *sneatnav.TUI) (sneatnav.Panel, error),
//) {
//	menu := NewDataTugMainMenu(tui, selectedMenuItem)
//
//	content, err := getContent(tui)
//	if err != nil {
//		panic(err)
//	}
//	if content == nil {
//		panic("getContent() returned nil")
//	}
//
//	// Allow keyboard navigation from the menu to the header with Shift+Tab (Backtab) or Up arrow.
//	// This enables Breadcrumbs to receive focus and thus its InputHandler to be called.
//
//	grid.SetFocusFunc(func() {
//		menu.TakeFocus()
//	})
//
//	_ = sneatnav.NewRow(tui.App,
//		menu,
//		content,
//	)
//}
