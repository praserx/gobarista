// Copyright 2023 PraserX
package main

import "github.com/urfave/cli/v2"

var FlagConfig = cli.StringFlag{
	Name:    "config",
	Value:   "./.local/config.ini",
	Usage:   "path to configuration file",
	Aliases: []string{"c"},
}
