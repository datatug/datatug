package commands

import (
	"context"

	"github.com/datatug/datatug-cli/apps/datatug"
	dtproject2 "github.com/datatug/datatug-cli/apps/datatug/datatugui/dtproject"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/dtsettings"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/dtviewers"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/dtviewers/clouds/aws/awsui"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/dtviewers/clouds/azure/azureui"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/dtviewers/clouds/gcloud/gcloudui"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/dtviewers/sqlviewer"
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

	tui := datatug.NewDatatugTUI()

	registerModules()

	if err := dtproject2.GoProjectsScreen(tui, sneatnav.FocusToMenu); err != nil {
		panic(err)
	}

	tui.App.SetRoot(tui.Grid, true)
	return tui.App.Run()
}

func registerModules() {

	dtproject2.RegisterModule()

	gcloudui.RegisterAsViewer()
	awsui.RegisterAsViewer()
	azureui.RegisterAsViewer()
	sqlviewer.RegisterAsViewer()

	dtviewers.RegisterModule()
	dtsettings.RegisterModule()
}
