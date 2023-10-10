// Copyright 2023 PraserX
package flags

import "github.com/urfave/cli/v2"

var FlagConfig = cli.StringFlag{
	Name:    "config",
	Value:   ".local/config.ini",
	Usage:   "path to configuration file",
	Aliases: []string{"c"},
}

var FlagPrettyPrint = cli.BoolFlag{
	Name:    "pretty",
	Aliases: []string{"p"},
	Usage:   "print formatted output",
}
