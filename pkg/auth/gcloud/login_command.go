package gcloud

import (
	"context"

	"github.com/urfave/cli/v3"
)

func loginCommand() *cli.Command {
	return &cli.Command{
		Name: "login",
		Action: func(ctx context.Context, command *cli.Command) error {
			return nil
		},
	}
}
