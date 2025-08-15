package commands

import (
	"context"
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-cli/apps/datatug/ui"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v3"
)

func uiCommandArgs() *cli.Command {
	return &cli.Command{
		Name:        "ui",
		Usage:       "Starts UI",
		Description: "",
		Action: func(ctx context.Context, c *cli.Command) error {
			v := &uiCommand{}
			return v.Execute(nil)
		},
	}
}

type uiCommand struct {
}

func (v *uiCommand) Execute(_ []string) error {

	app := tview.NewApplication()
	app.EnableMouse(true)
	//app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	//	switch event.Key() {
	//	case tcell.KeyTab:
	//		// Move to next (default behavior)
	//		return event
	//	case tcell.KeyBacktab: // This is Shift+Tab
	//		// Move to previous
	//		app.SetFocus(getPreviousFocusable())
	//		return nil // Consume the event
	//	}
	//	return event
	//
	//})
	tui := tapp.NewTUI(app)
	if err := ui.GoHomeScreen(tui); err != nil {
		panic(err)
	}
	app.SetRoot(tui.Grid, true)
	return app.Run()
}
