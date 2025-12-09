package commands

import (
	"context"

	"github.com/datatug/datatug-cli/apps/datatug"
	dtprojects2 "github.com/datatug/datatug-cli/apps/datatug/datatugui/dtscreeens/dtprojects"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/dtscreeens/dtsettings"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers/clouds/aws/awsui"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers/clouds/azure/azureui"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers/clouds/gcloud/gcloudui"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers/sqlviewer"
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

	if err := dtprojects2.GoProjectsScreen(tui, sneatnav.FocusToMenu); err != nil {
		panic(err)
	}

	tui.App.SetRoot(tui.Grid, true)
	return tui.App.Run()
}

func registerModules() {

	dtprojects2.RegisterModule()

	gcloudui.RegisterAsViewer()
	awsui.RegisterAsViewer()
	azureui.RegisterAsViewer()
	sqlviewer.RegisterAsViewer()

	viewers.RegisterModule()
	dtsettings.RegisterModule()
}
