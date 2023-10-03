// Copyright 2023 PraserX
package helpers

import (
	"github.com/praserx/gobarista/pkg/config"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/mail"
	"github.com/urfave/cli/v2"
	"gopkg.in/ini.v1"
)

func GetSMTPConfig(cfg *ini.File) *mail.MailOptions {
	return &mail.MailOptions{
		Host: cfg.Section("smtp").Key("host").String(),
		Port: cfg.Section("smtp").Key("port").String(),
		User: cfg.Section("smtp").Key("username").String(),
		Pass: cfg.Section("smtp").Key("password").String(),
		Name: cfg.Section("smtp").Key("name").String(),
		From: cfg.Section("smtp").Key("from").String(),
	}
}

func SetupMail(ctx *cli.Context) (err error) {
	if err = config.Load(ctx.String("config")); err != nil {
		return err
	}

	opts := GetSMTPConfig(config.Get())
	options := []mail.Option{}
	options = append(options, mail.WithHost(opts.Host))
	options = append(options, mail.WithPort(opts.Port))
	options = append(options, mail.WithName(opts.Name))
	options = append(options, mail.WithFrom(opts.From))
	options = append(options, mail.WithCredentials(opts.User, opts.Pass))
	mail.SetupMailer(options...)

	return nil
}

func GetDatabaseConfig(cfg *ini.File) *database.DatabaseOptions {
	return &database.DatabaseOptions{
		Path: cfg.Section("paths").Key("database").String(),
	}
}

func SetupDatabase(ctx *cli.Context) (err error) {
	if err = config.Load(ctx.String("config")); err != nil {
		return err
	}

	opts := GetDatabaseConfig(config.Get())
	options := []database.Option{}
	options = append(options, database.WithPath(opts.Path))
	database.SetupDatabase(options...)

	return nil
}
