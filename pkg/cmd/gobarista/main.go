// Copyright 2023 PraserX
package main

import (
	"os"

	"github.com/praserx/gobarista/pkg/cmd/gobarista/commands"
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
			&FlagConfig,
		},
		Commands: []*cli.Command{
			&commands.Billing,
			&commands.Database,
			&commands.Users,
		},
	}

	if err = app.Run(os.Args); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// func main() {
// 	var err error

// 	app := &cli.App{
// 		Name:        "GoBarista",
// 		Usage:       "Coffee for everyone, easy bill for our coffee samaritan.",
// 		UsageText:   "gobarista config_file.ini\ngobarista [global options] command [command options] [arguments...]",
// 		Version:     version.VERSION,
// 		Description: "This is a simple coffee billing application - for enterprise, non-profits or government.",
// 		Copyright:   "(c) Praser",
// 		Commands: []*cli.Command{
// 			&commandTestMail,
// 		},
// 		Action: func(c *cli.Context) error {
// 			if c.NArg() > 0 {
// 				var err error

// 				if err = config.Load(c.Args().Get(0)); err != nil {
// 					return err
// 				}

// 				if config.Get().Section("").Key("app_mode").String() == "development" {
// 					logger.SetupLogger([]logger.Option{logger.WithDebugMode(true)}...)
// 				}

// 				setupMail(config.Get())
// 				setupDatabase(config.Get())

// 				var serverOptions []web.Option
// 				serverOptions = append(serverOptions, web.WithHost(config.Get().Section("server").Key("host").String()))
// 				serverOptions = append(serverOptions, web.WithPort(config.Get().Section("server").Key("port").String()))
// 				serverOptions = append(serverOptions, web.WithBaseUrl(config.Get().Section("server").Key("base_url").String()))
// 				if config.Get().Section("").Key("app_mode").String() == "development" {
// 					serverOptions = append(serverOptions, web.WithDebug(true))
// 				}

// 				var ws *web.Server
// 				if ws, err = web.NewServer(serverOptions...); err != nil {
// 					return err
// 				}

// 				if err = ws.Serve(); err != nil {
// 					return err
// 				}

// 				return nil
// 			}

// 			return fmt.Errorf("missing config file")
// 		},
// 	}

// 	if err = app.Run(os.Args); err != nil {
// 		fmt.Fprintf(os.Stderr, "fatal error: %s\n", err.Error())
// 		os.Exit(1)
// 	}
// }
