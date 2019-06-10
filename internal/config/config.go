// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package config

import (
	"encoding/json"
	"expvar"
	"fmt"
	"io"

	"github.com/circonus-labs/circonus-cloud-agent/internal/config/defaults"
	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// Log defines the running config.log structure
type Log struct {
	Level  string `json:"level" yaml:"level" toml:"level"`
	Pretty bool   `json:"pretty" yaml:"pretty" toml:"pretty"`
}

// // API defines the running config.api structure
// type API struct {
// 	App    string `json:"app" yaml:"app" toml:"app"`
// 	CAFile string `mapstructure:"ca_file" json:"ca_file" yaml:"ca_file" toml:"ca_file"`
// 	Key    string `json:"key" yaml:"key" toml:"key"`
// 	URL    string `json:"url" yaml:"url" toml:"url"`
// 	Debug  bool   `json:"debug" yaml:"debug" toml:"debug"`
// }

// Config defines the running config structure
type Config struct {
	// API   API          `json:"api" yaml:"api" toml:"api"`
	Debug       bool         `json:"debug" yaml:"debug" toml:"debug"`
	Log         Log          `json:"log" yaml:"log" toml:"log"`
	AWS         *AWSConfig   `json:"aws" toml:"aws" yaml:"aws"`
	Azure       *AzureConfig `json:"azure" toml:"azure" yaml:"azure"`
	GCP         *GCPConfig   `json:"gcp" toml:"gcp" yaml:"gcp"`
	PipeSubmits bool         `json:"pipe_submits" toml:"pipe_submits" yaml:"pipe_submits"`
}

// AWSConfig defines the AWS cloud service configuration
type AWSConfig struct {
	Enabled bool   `json:"enabled" toml:"enabled" yaml:"enabled"`
	ConfDir string `json:"conf_dir" toml:"conf_dir" yaml:"conf_dir"`
}

// AzureConfig defines the Azure cloud service configuration
type AzureConfig struct {
	Enabled bool   `json:"enabled" toml:"enabled" yaml:"enabled"`
	ConfDir string `json:"conf_dir" toml:"conf_dir" yaml:"conf_dir"`
}

// GCPConfig defines the AWS cloud service configuration
type GCPConfig struct {
	Enabled bool   `json:"enabled" toml:"enabled" yaml:"enabled"`
	ConfDir string `json:"conf_dir" toml:"conf_dir" yaml:"conf_dir"`
}

//
// NOTE: adding a Key* MUST be reflected in the Config structures above
//
const (
	// // KeyAPICAFile custom ca for circonus api (e.g. inside)
	// KeyAPICAFile = "api.ca_file"
	//
	// // KeyAPITokenApp circonus api token key application name
	// KeyAPITokenApp = "api.app"
	//
	// // KeyAPITokenKey circonus api token key
	// KeyAPITokenKey = "api.key"
	//
	// // KeyAPIURL custom circonus api url (e.g. inside)
	// KeyAPIURL = "api.url"
	//
	// // KeyAPIDebug turns on debugging for circonus api calls
	// KeyAPIDebug = "api.debug"

	// KeyDebug enables debug messages
	KeyDebug = "debug"

	// KeyLogLevel logging level (panic, fatal, error, warn, info, debug, disabled)
	KeyLogLevel = "log.level"

	// KeyLogPretty output formatted log lines (for running in foreground)
	KeyLogPretty = "log.pretty"

	// KeyShowConfig - show configuration and exit
	KeyShowConfig = "show-config"

	// KeyShowVersion - show version information and exit
	KeyShowVersion = "version"

	// KeyPipeSubmits - use io pipe for metric submissions (experimental)
	// KeyPipeSubmits = "pipe_submits"
)

var (
	// MetricNameSeparator defines character used to delimit metric name parts
	MetricNameSeparator = defaults.MetricNameSeparator // var, TBD whether it will become configurable
)

// Validate verifies the required portions of the configuration
func Validate() error {

	// err := validateAPIOptions()
	// if err != nil {
	// 	return errors.Wrap(err, "API config")
	// }

	return nil
}

// StatConfig adds the running config to the app stats
func StatConfig() error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	// cfg.API.Key = "..."
	// cfg.API.App = "..."

	expvar.Publish("config", expvar.Func(func() interface{} {
		return &cfg
	}))

	return nil
}

// getConfig dumps the current configuration and returns it
func getConfig() (*Config, error) {
	var cfg *Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "parsing config")
	}

	return cfg, nil
}

// ShowConfig prints the running configuration
func ShowConfig(w io.Writer) error {
	var cfg *Config
	var err error
	var data []byte

	cfg, err = getConfig()
	if err != nil {
		return err
	}

	format := viper.GetString(KeyShowConfig)

	log.Debug().Str("format", format).Msg("show-config")

	switch format {
	case "json":
		data, err = json.MarshalIndent(cfg, " ", "  ")
		if err != nil {
			return errors.Wrap(err, "formatting config (json)")
		}
	case "yaml":
		data, err = yaml.Marshal(cfg)
		if err != nil {
			return errors.Wrap(err, "formatting config (yaml)")
		}
	case "toml":
		data, err = toml.Marshal(*cfg)
		if err != nil {
			return errors.Wrap(err, "formatting config (toml)")
		}
	default:
		return errors.Errorf("unknown config format '%s'", format)
	}

	_, err = fmt.Fprintf(w, "\n%s\n", data)
	return err
}
