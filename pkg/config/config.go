// Copyright 2023 PraserX
package config

import (
	"gopkg.in/ini.v1"
)

// Global logger instance
var Configuration *ini.File

// Load function loads ini configuration from file to global Configuration struct.
func Load(path string) (err error) {
	if Configuration, err = ini.Load(path); err != nil {
		return err
	}
	return nil
}

// Ge returns configuration object.
func Get() *ini.File {
	return Configuration
}
