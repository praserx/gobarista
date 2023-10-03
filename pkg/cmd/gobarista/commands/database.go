// Copyright 2023 PraserX
package commands

import (
	"github.com/praserx/gobarista/pkg/cmd/gobarista/helpers"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/urfave/cli/v2"
)

var Database = cli.Command{
	Name:  "database",
	Usage: "Database initialization and checks",
	Subcommands: []*cli.Command{
		&DatabaseInitialize,
		&DatabaseMigrate,
	},
}

var DatabaseInitialize = cli.Command{
	Name:  "initialize",
	Usage: "Initialize SQLite database",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}
		database.RunAutoMigration()
		return nil
	},
}

var DatabaseMigrate = cli.Command{
	Name:  "migrate",
	Usage: "Migrate SQLite database",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}
		database.RunAutoMigration()
		return nil
	},
}

var DatabaseCheck = cli.Command{
	Name:  "check",
	Usage: "Check SQLite database",
	Action: func(ctx *cli.Context) (err error) {
		return helpers.SetupDatabase(ctx)
	},
}
