// Copyright 2023 PraserX
package database

// Options are used for Database construct function.
type DatabaseOptions struct {
	Path string
}

// Option definition
type Option func(*DatabaseOptions)

// WithPath option specification.
func WithPath(path string) Option {
	return func(opts *DatabaseOptions) {
		opts.Path = path
	}
}
