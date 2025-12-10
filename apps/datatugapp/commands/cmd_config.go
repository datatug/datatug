package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/urfave/cli/v3"
)

func configCommandAction(_ context.Context, _ *cli.Command) error {
	settings, err := appconfig.GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}
	if err = appconfig.PrintSettings(settings, appconfig.FormatYaml, os.Stdout); err != nil {
		return err
	}
	return nil
}

func configCommand() *cli.Command {
	return &cli.Command{
		Name:        "config",
		Usage:       "Prints config",
		Description: "",
		Action:      configCommandAction,
	}
}
