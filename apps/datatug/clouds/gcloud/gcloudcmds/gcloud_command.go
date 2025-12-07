package gcloudcmds

import "github.com/urfave/cli/v3"

func GoogleCloudCommand() *cli.Command {
	return &cli.Command{
		Name: "gcloud",
		Commands: []*cli.Command{
			loginCommand(),
			projectsCommand(),
		},
	}
}
