package commands

import (
	"context"

	"github.com/datatug/datatug-cli/apps/datatug/dtscreeens/clouds"
	"github.com/datatug/datatug-cli/apps/datatug/dtscreeens/dthome"
	"github.com/datatug/datatug-cli/apps/datatug/dtscreeens/dtprojects"
	"github.com/datatug/datatug-cli/apps/datatug/dtscreeens/dtsettings"
	"github.com/datatug/datatug-cli/apps/datatug/dtscreeens/dtviewers"
	"github.com/datatug/datatug-cli/apps/firestoreviewer/fsviewer"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
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

	var tui *sneatnav.TUI
	tui = sneatnav.NewTUI(app, sneatv.NewBreadcrumb(" ⛴ DataTug", func() error {
		return dthome.GoHomeScreen(tui, sneatnav.FocusToContent)
	}))

	registerModules(tui)

	if err := dthome.GoHomeScreen(tui, sneatnav.FocusToContent); err != nil {
		panic(err)
	}

	app.SetRoot(tui.Grid, true)
	return app.Run()
}

func registerModules(tui *sneatnav.TUI) {

	// Main menu screens
	dthome.RegisterModule()
	clouds.RegisterModule()
	dtsettings.RegisterModule()
	dtprojects.RegisterModule()
	dtviewers.RegisterModule()

	// Sub-modules
	fsviewer.RegisterModule(tui)
}
