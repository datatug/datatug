package commands

import (
	"context"

	"github.com/datatug/datatug-cli/apps/datatug/clouds"
	"github.com/datatug/datatug-cli/apps/datatug/clouds/gcloud/gcloudui"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
	"github.com/datatug/datatug-cli/apps/datatug/dtscreeens/dtprojects"
	"github.com/datatug/datatug-cli/apps/datatug/dtscreeens/dtsettings"
	"github.com/datatug/datatug-cli/apps/datatug/dtscreeens/dtviewers"
	"github.com/datatug/datatug-cli/apps/firestoreviewer/fsviewer"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
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

	tui := datatugui.NewDatatugTUI()

	registerModules(tui)

	if err := datatugui.GoHomeScreen(tui, sneatnav.FocusToContent); err != nil {
		panic(err)
	}

	tui.App.SetRoot(tui.Grid, true)
	return tui.App.Run()
}

func registerModules(tui *sneatnav.TUI) {

	// Main menu screens
	datatugui.RegisterModule()
	clouds.RegisterModule([]clouds.Cloud{
		{
			Name:     "Google Cloud",
			Shortcut: 'g',
			Action: func(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
				return gcloudui.GoHome(&gcloudui.GCloudContext{TUI: tui}, focusTo)
			},
		},
		{
			Name:     "Amazon Web Services",
			Shortcut: 'a',
		},
		{
			Name:     "Microsoft Azure",
			Shortcut: 'm',
		},
	})
	dtsettings.RegisterModule()
	dtprojects.RegisterModule()
	dtviewers.RegisterModule()

	// Sub-modules
	fsviewer.RegisterModule(tui)
}
