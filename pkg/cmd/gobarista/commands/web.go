// Copyright 2023 PraserX
package commands

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/praserx/gobarista/pkg/cmd/gobarista/flags"
	"github.com/praserx/gobarista/pkg/cmd/gobarista/helpers"
	"github.com/praserx/gobarista/pkg/config"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/webserver"
	"github.com/urfave/cli/v2"
)

var Web = cli.Command{
	Name:    "web",
	Aliases: []string{"w"},
	Usage:   "Application web server operator",
	Flags: []cli.Flag{
		&flags.FlagPlainPrint,
	},
	Subcommands: []*cli.Command{
		&WebRun,
	},
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		users, err := database.SelectAllUsers()
		if err != nil {
			return fmt.Errorf("error: cannot get users: %v", err.Error())
		}

		if !ctx.Bool("plain") {
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

var WebRun = cli.Command{
	Name:  "run",
	Usage: "Start web server application",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if err = config.Load(ctx.String("config")); err != nil {
			return err
		}

		host := config.Get().Section("server").Key("host").String()
		port, _ := config.Get().Section("server").Key("port").Int()

		srv := webserver.NewServer(host, port)
		srv.Run()

		return nil
	},
}
