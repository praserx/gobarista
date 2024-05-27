// Copyright 2023 PraserX
package main

import (
	"os"

	"github.com/praserx/gobarista/pkg/cmd/gobarista/commands"
	"github.com/praserx/gobarista/pkg/cmd/gobarista/flags"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/version"

	"github.com/urfave/cli/v2"
)

func main() {
	var err error

	app := &cli.App{
		Name:        "gobarista",
		Usage:       "Coffee for everyone, easy bill for our coffee samaritan.",
		UsageText:   "gobarista -c config_file.ini\ngobarista [global options] command [command options] [arguments...]",
		Version:     version.VERSION,
		Description: "This is a simple coffee billing application - for enterprise, non-profits or government.",
		Copyright:   "(c) Praser",
		Flags: []cli.Flag{
			&flags.FlagConfig,
		},
		Commands: []*cli.Command{
			&commands.Billing,
			&commands.Database,
			&commands.Users,
			&commands.Web,
		},
	}

	if err = app.Run(os.Args); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
