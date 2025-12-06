package commands

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/datatug/datatug-cli/apps"
	"github.com/datatug/datatug-cli/apps/datatug/dtbubble"
	"github.com/datatug/datatug-cli/apps/firestoreviewer"
	"github.com/datatug/datatug-cli/pkg/auth"
	"github.com/datatug/datatug-cli/pkg/auth/gcloud"
	"github.com/urfave/cli/v3"
)

func DatatugCommand() *cli.Command {
	return &cli.Command{
		Action: datatugCommandAction,
		Flags:  []cli.Flag{apps.TUIFlag},
		Commands: []*cli.Command{
			initCommand(),
			auth.AuthCommand(),
			gcloud.GoogleCloudCommand(),
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
			firestoreviewer.FirestoreCommand(),
		},
	}
}

func datatugCommandAction(_ context.Context, cmd *cli.Command) error {
	if !apps.TUIFlag.IsSet() {
		// Show default help text when TUI is not requested
		_ = cli.ShowRootCommandHelp(cmd)
		return nil
	}
	datatugApp := dtbubble.DatatugAppModel()
	p := tea.NewProgram(datatugApp, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		// Ensure the error is printed to the console explicitly
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return nil
}
