package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func getCLIApp() *cli.App {
	app := cli.NewApp()
	app.Name = Name
	app.HelpName = Name
	app.Usage = fmt.Sprintf("%s [command] [arguments]", Name)
	app.Version = Name
	app.Flags = []cli.Flag{
		cli.VersionFlag,
		cli.HelpFlag,
		cli.StringFlag{
			Name:  "config, c",
			Usage: "config file path, default(/etc/kubetrack/config.yaml)",
			Value: "/etc/kubetrack/config.yaml",
		},
	}
	app.Action = func(c *cli.Context) error {
		return runMain(c)
	}
	return app
}
