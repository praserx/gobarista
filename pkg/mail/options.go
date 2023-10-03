// Copyright 2023 PraserX
package mail

type MailOptions struct {
	Host  string
	Port  string
	Name  string
	From  string
	User  string
	Pass  string
	Links map[string]string
}

type Option func(*MailOptions)

func WithHost(host string) Option {
	return func(opts *MailOptions) {
		opts.Host = host
	}
}

func WithPort(port string) Option {
	return func(opts *MailOptions) {
		opts.Port = port
	}
}

func WithName(name string) Option {
	return func(opts *MailOptions) {
		opts.Name = name
	}
}

func WithFrom(from string) Option {
	return func(opts *MailOptions) {
		opts.From = from
	}
}

func WithCredentials(username, password string) Option {
	return func(opts *MailOptions) {
		opts.User = username
		opts.Pass = password
	}
}

func WithLinks(links map[string]string) Option {
	return func(opts *MailOptions) {
		opts.Links = links
	}
}
