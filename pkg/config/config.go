// Package config provides utilities for loading structured configuration from
// files into typed Go structs. It supports JSON, TOML, YAML, and environment
// variable overrides via the configor library.
//
// The config directory is resolved from the CONFIG_PATH environment variable;
// if that variable is unset, it falls back to the directory of this source
// file so the package works out-of-the-box without any runtime configuration.
package config

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/jinzhu/configor"
	"github.com/samber/lo"
)

// LoadConfig loads one or more config files into a value of type T and returns
// a pointer to the populated struct. File paths are relative to the config
// directory (see getConfigLocation). Multiple files are merged in order, with
// later files taking precedence. Environment variables can override any field
// according to configor conventions.
//
// Example:
//
//	type AppConfig struct {
//	    Port int    `json:"port"`
//	    DSN  string `json:"dsn"`
//	}
//
//	cfg, err := config.LoadConfig[AppConfig]("app.yaml")
func LoadConfig[T any](files ...string) (*T, error) {
	var result T
	cfgPath := getConfigLocation()
	filesToLoad := lo.Map(files, func(item string, index int) string {
		return fmt.Sprintf("%s/%s", cfgPath, item)
	})
	err := configor.Load(&result, filesToLoad...)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return &result, nil
}

// getConfigLocation returns the directory from which config files are loaded.
// It checks the CONFIG_PATH environment variable first; if unset it falls back
// to the directory that contains this source file, which makes the package
// self-contained when config files are bundled alongside the binary source.
func getConfigLocation() string {
	configPath, exists := os.LookupEnv("CONFIG_PATH")
	if exists {
		return configPath
	}

	_, filename, _, _ := runtime.Caller(0) //nolint:dogsled // we actually want to ignore the rest

	return path.Join(path.Dir(filename), "./")
}
