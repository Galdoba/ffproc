package main

import (
	"github.com/Galdoba/ffproc/cmd/procman/commands"
	"github.com/urfave/cli/v2"
)

func StartProcman() *cli.App {
	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = programName
	app.Usage = "manager for processing media files"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Category:    "",
			DefaultText: co,
			Usage:       "non-default config file",

			Value:       "",
			Destination: new(string),
		},
	}

	app.Before = func(c *cli.Context) error {
		return nil
	}
	app.Commands = []*cli.Command{
		commands.Assemble(),
	}
	// app.DefaultCommand = "run"

	app.After = func(c *cli.Context) error {
		return nil
	}
	return app
}
