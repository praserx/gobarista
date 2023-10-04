// Copyright 2023 PraserX
package commands

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/praserx/gobarista/pkg/cmd/gobarista/helpers"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/models"
	"github.com/urfave/cli/v2"
)

var Users = cli.Command{
	Name:  "users",
	Usage: "User management operations",
	Subcommands: []*cli.Command{
		&UsersAdd,
		&UsersAddBulk,
	},
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		users, err := database.SelectAllUsers()
		if err != nil {
			return fmt.Errorf("error: cannot get users: %v", err.Error())
		}

		for _, user := range users {
			fmt.Printf("%s %s\t\t%s\t\t%s\n", user.Firstname, user.Lastname, user.Email, user.Location)
		}

		return nil
	},
}

var UsersAdd = cli.Command{
	Name:      "add",
	Usage:     "Add new user",
	ArgsUsage: "[employee_id firstname lastname e-mail location]",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if ctx.NArg() != 5 {
			return fmt.Errorf("error: too few arguments: requires (5), get (%d)", ctx.NArg())
		}

		user := models.User{
			EID:       ctx.Args().Get(0),
			Firstname: ctx.Args().Get(1),
			Lastname:  ctx.Args().Get(2),
			Email:     ctx.Args().Get(3),
			Location:  ctx.Args().Get(4),
		}

		if _, err = database.SelectUserByEID(user.EID); err != nil {
			id, err := database.InsertUser(user)
			if err != nil {
				return fmt.Errorf("error: cannot create user: %v", err.Error())
			}
			logger.Info(fmt.Sprintf("user successfully created: new user id: %d", id))
		} else {
			logger.Info("user already exists")
		}

		return nil
	},
}

var UsersAddBulk = cli.Command{
	Name:      "add-bulk",
	Usage:     "Add multiple users at once by csv file (name,e-mail,location,...)",
	ArgsUsage: "[file.csv]",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if ctx.NArg() != 1 {
			fmt.Println("Too few arguments")
		}

		fmt.Println(ctx.Args().Get(0))

		const (
			CSV_EID = iota
			CSV_NAME
			CSV_EMAIL
			CSV_LOCATION
			CSV_COFFEES
		)

		var fd *os.File
		if fd, err = os.Open(ctx.Args().Get(0)); err != nil {
			return err
		}

		ioReader := bufio.NewReader(fd)
		reader := csv.NewReader(ioReader)

		for {
			var record []string
			if record, err = reader.Read(); err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			user := models.User{
				EID:       record[CSV_EID],
				Firstname: strings.Fields(record[CSV_NAME])[0],
				Lastname:  strings.Fields(record[CSV_NAME])[1],
				Email:     record[CSV_EMAIL],
				Location:  record[CSV_LOCATION],
			}

			if _, err = database.SelectUserByEID(user.EID); err != nil {
				id, err := database.InsertUser(user)
				if err != nil {
					return fmt.Errorf("error: cannot create user: %v", err.Error())
				}
				logger.Info(fmt.Sprintf("user successfully created: new user id: %d", id))
			} else {
				logger.Info("user already exists")
			}
		}

		return nil
	},
}