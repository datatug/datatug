package commands

import (
	"context"

	"github.com/datatug/datatug-cli/apps/datatug/clouds"
	"github.com/datatug/datatug-cli/apps/datatug/clouds/aws/awsui"
	"github.com/datatug/datatug-cli/apps/datatug/clouds/azure/azureui"
	"github.com/datatug/datatug-cli/apps/datatug/clouds/gcloud/gcloudui"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
	"github.com/datatug/datatug-cli/apps/datatug/dtscreeens/dtprojects"
	"github.com/datatug/datatug-cli/apps/datatug/dtscreeens/dtsettings"
	"github.com/datatug/datatug-cli/apps/datatug/dtscreeens/dtviewers"
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

	dtprojects.RegisterModule()
	// Main menu screens
	datatugui.RegisterModule()

	cloudContext := &clouds.CloudContext{TUI: tui}
	clouds.RegisterModule([]clouds.Cloud{
		{
			Name:     "Google Cloud",
			Shortcut: 'g',
			Action: func(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
				return gcloudui.GoHome(&gcloudui.GCloudContext{
					CloudContext: cloudContext,
				}, focusTo)
			},
		},
		{
			Name:     "Amazon Web Services",
			Shortcut: 'a',
			Action: func(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
				return awsui.GoAwsHome(cloudContext, focusTo)
			},
		},
		{
			Name:     "Microsoft Azure",
			Shortcut: 'm',
			Action: func(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
				return azureui.GoAzureHome(cloudContext, focusTo)
			},
		},
	})
	dtsettings.RegisterModule()
	dtviewers.RegisterModule()
}
