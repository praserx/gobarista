// Copyright 2023 PraserX
package commands

import (
	"fmt"

	"github.com/praserx/gobarista/pkg/cmd/gobarista/helpers"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/models"
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

		_, err = database.InsertSchema(models.Schema{Version: models.VERSION})
		if err != nil {
			return fmt.Errorf("cannot update schema version for database: %v", err)
		}

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

		schema, err := database.SelectVersion()
		if err != nil {
			return fmt.Errorf("cannot get schema from database for version check: %v", err)
		}

		if schema.Version != models.VERSION {
			err = database.UpdateVersion(models.VERSION)
			if err != nil {
				return fmt.Errorf("cannot update schema version for database: %v", err)
			}
		}

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
