// Copyright 2023 PraserX
package commands

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/praserx/gobarista/pkg/cmd/gobarista/flags"
	"github.com/praserx/gobarista/pkg/cmd/gobarista/helpers"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/models"
	"github.com/urfave/cli/v2"
)

var Users = cli.Command{
	Name:    "users",
	Aliases: []string{"u"},
	Usage:   "User management operations",
	Flags: []cli.Flag{
		&flags.FlagPrettyPrint,
	},
	Subcommands: []*cli.Command{
		&UsersAdd,
		&UsersAddBulk,
		&UsersContacts,
	},
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		users, err := database.SelectAllUsers()
		if err != nil {
			return fmt.Errorf("error: cannot get users: %v", err.Error())
		}

		if ctx.Bool("pretty") {
			t := table.NewWriter()
			t.AppendHeader(table.Row{"ID", "Name", "E-mail", "Location"})
			for _, user := range users {
				t.AppendRow(table.Row{user.ID, user.Firstname + " " + user.Lastname, user.Email, user.Location})
			}
			fmt.Println(t.Render())
		} else {
			for _, user := range users {
				fmt.Printf("%d %s %s, %s, %s\n", user.ID, user.Firstname, user.Lastname, user.Email, user.Location)
			}
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

var UsersContacts = cli.Command{
	Name:    "contacts",
	Aliases: []string{"c"},
	Usage:   "Get e-mail contacts of all customers",
	Flags: []cli.Flag{
		&flags.FlagPrettyPrint,
	},
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if ctx.NArg() != 0 {
			return fmt.Errorf("error: too few arguments: requires (1), get (%d)", ctx.NArg())
		}

		bills, err := database.SelectAllBills()
		if err != nil {
			return fmt.Errorf("error: cannot get billing period: %v", err)
		}

		var contacts []string
		for _, bill := range bills {
			user, err := database.SelectUserByID(bill.UserID)
			if err != nil {
				logger.Error(fmt.Sprintf("error: cannot get user by id: user_id=%d: %v", bill.UserID, err.Error()))
			}

			if !slices.Contains(contacts, user.Email) {
				contacts = append(contacts, user.Email)
				fmt.Println(user.Email)
			}
		}

		return nil
	},
}
