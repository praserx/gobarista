// Copyright 2023 PraserX
package mail

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/resources"
	"github.com/praserx/mailgo"
)

var mail *Mailer

// ErrNilMailer specifies nil mailer instance error
var ErrNilMailer = fmt.Errorf("fatal: mailer instance is nil")

// Mailer struct definition.
type Mailer struct {
	mailer *mailgo.Mailer
}

// SetupMailer set ups mailer.
func SetupMailer(opts ...Option) (err error) {
	mail, err = NewMailer(opts...)
	return err
}

// CheckMail checks if mailer instance exists.
func CheckMail() error {
	if mail == nil {
		return ErrNilMailer
	}
	return nil
}

// Mailer returns mailer instance.
func Mail() *Mailer {
	return mail
}

// NewMailer creates new mailer instance. If we can't create new mailer than
// error is returned.
func NewMailer(opts ...Option) (mailer *Mailer, err error) {
	var options = &MailOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var mailgoOptions []mailgo.MailerOption
	mailgoOptions = append(mailgoOptions, mailgo.WithHost(options.Host))
	mailgoOptions = append(mailgoOptions, mailgo.WithPort(options.Port))
	mailgoOptions = append(mailgoOptions, mailgo.WithFrom(options.From))
	mailgoOptions = append(mailgoOptions, mailgo.WithName(options.Name))
	if options.User != "" && options.Pass != "" {
		mailgoOptions = append(mailgoOptions, mailgo.WithCredentials(options.User, options.Pass))
	}

	var mailgoMailer *mailgo.Mailer
	if mailgoMailer, err = mailgo.NewMailer(mailgoOptions...); err != nil {
		return nil, fmt.Errorf("cannot initialize mailer: mailgo: %v", err)
	}

	mailer = &Mailer{}
	mailer.mailer = mailgoMailer

	return mailer, nil
}

// SendMail ...
func SendMail(recipient string, es *EmailSettings) (err error) {
	return mail.mailer.SendMail([]string{recipient}, es.Subject, es.Plain, es.HTML)
}

// GetBillHTMLTemplate ...
func GetBillHTMLTemplate(vars BillTemplateVars) string {
	var err error
	var tmpl *template.Template
	var buffer bytes.Buffer

	if tmpl, err = template.ParseFS(resources.DirTemplates, resources.HTML_BILL_TEMPLATE_FULL_PATH); err != nil {
		logger.Error("cannot parse template: " + err.Error())
		return ""
	}

	if err = tmpl.Execute(&buffer, vars); err != nil {
		logger.Error("cannot parse template: " + err.Error())
		return ""
	}

	return buffer.String()
}

// GetBillHTMLTemplate ...
func GetConfirmHTMLTemplate(vars ConfirmationTemplateVars) string {
	var err error
	var tmpl *template.Template
	var buffer bytes.Buffer

	if tmpl, err = template.ParseFS(resources.DirTemplates, resources.HTML_CONFIRM_TEMPLATE_FULL_PATH); err != nil {
		logger.Error("cannot parse template: " + err.Error())
		return ""
	}

	if err = tmpl.Execute(&buffer, vars); err != nil {
		logger.Error("cannot parse template: " + err.Error())
		return ""
	}

	return buffer.String()
}

// GetUnpaidHTMLTemplate ...
func GetUnpaidHTMLTemplate(vars UnpaidTemplateVars) string {
	var err error
	var tmpl *template.Template
	var buffer bytes.Buffer

	if tmpl, err = template.ParseFS(resources.DirTemplates, resources.HTML_UNPAID_TEMPLATE_FULL_PATH); err != nil {
		logger.Error("cannot parse template: " + err.Error())
		return ""
	}

	if err = tmpl.Execute(&buffer, vars); err != nil {
		logger.Error("cannot parse template: " + err.Error())
		return ""
	}

	return buffer.String()
}
