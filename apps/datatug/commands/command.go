package commands

import (
	"context"

	"github.com/datatug/datatug-cli/apps"
	"github.com/datatug/datatug-cli/apps/datatug/datatugui/viewers/clouds/gcloud/gcloudcmds"
	"github.com/datatug/datatug-cli/pkg/auth"
	"github.com/urfave/cli/v3"
)

func DatatugCommand() *cli.Command {
	return &cli.Command{
		Action: datatugCommandAction,
		Flags:  []cli.Flag{apps.TUIFlag},
		Commands: []*cli.Command{
			initCommand(),
			auth.AuthCommand(),
			gcloudcmds.GoogleCloudCommand(),
			configCommand(),
			datasetCommands(),
			datasetDefCommandArgs(),
			datasetDataCommandArgs(),
			datasetsCommandArgs(),
			demoCommandArgs(),
			updateUrlConfigCommandArgs(),
			projectsCommandArgs(),
			queriesCommand(),
			renderCommandArgs(),
			scanCommandArgs(),
			serveCommandArgs(),
			showCommandArgs(),
			uiCommandArgs(),
			testCommandArgs(),
			consoleCommandArgs(),
		},
	}
}

func datatugCommandAction(_ context.Context, cmd *cli.Command) error {
	if !apps.TUIFlag.IsSet() {
		// Show default help text when TUI is not requested
		_ = cli.ShowRootCommandHelp(cmd)
		return nil
	}
	return nil
}
