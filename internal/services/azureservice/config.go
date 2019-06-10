// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package azureservice

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/circonus-labs/circonus-cloud-agent/internal/release"
	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Config for instance of Azure metric collection service
type Config struct {
	ID       string                 `json:"id" toml:"id" yaml:"id"`                   // unique id for this service client instance, no spaces (ties several things together, short and immutable - logging, check search/create, tags, etc.)
	Azure    AzureConfig            `json:"azure" toml:"azure" yaml:"azure"`          // REQUIRED, azure configuration
	Circonus circonus.ServiceConfig `json:"circonus" toml:"circonus" yaml:"circonus"` // REQUIRED, circonus config: api credentials, check, broker, etc.
	Tags     circonus.Tags          `json:"tags" toml:"tags" yaml:"tags"`             // global tags, added to all metrics
}

// AzureConfig defines the Azure sdk credentials
type AzureConfig struct {
	DirectoryID       string `json:"directory_id" toml:"directory_id" yaml:"directory_id"`                   // aka tenant id
	ApplicationID     string `json:"application_id" toml:"application_id" yaml:"application_id"`             // aka client id
	ApplicationSecret string `json:"application_secret" toml:"application_secret" yaml:"application_secret"` // aka client secret
	SubscriptionID    string `json:"subscription_id" toml:"subscription_id" yaml:"subscription_id"`          // azure subscription id
	ResourceFilter    string `json:"resource_filter" toml:"resource_filter" yaml:"resource_filter"`          // filter expression for which resources to collect metrics about -- limit resources from which metrics are collected, otherwise **ALL** metrics from **ALL** resources will be collected. Suggested method, add a tag to each resource from which to collect metrics (e.g. Tag Name:circonus and Tag Value:enabled) then use the expression: `tagName eq 'circonus' and tagValue eq 'enabled'`
	CloudName         string `json:"cloud_name" toml:"cloud_name" yaml:"cloud_name"`                         // DEFAULT: "AzurePublicCloud"
	UserAgent         string `json:"user_agent" toml:"user_agent" yaml:"user_agent"`                         // DEFAULT: "circonus-cloud-agent"
	Interval          int    `json:"collect_interval" toml:"collect_interval" yaml:"collect_interval"`       // DEFAULT: 5 - How often to collect metrics.(minutes >= 5)
}

const (
	defaultCloudName = "AzurePublicCloud"
	defaultInterval  = 5
)

func showExampleConfig(format string, w io.Writer) error {
	var err error
	var data []byte

	cfg := defaultConfig()

	switch format {
	case "json":
		data, err = json.MarshalIndent(cfg, " ", "  ")
	case "yaml":
		data, err = yaml.Marshal(cfg)
	case "toml":
		data, err = toml.Marshal(*cfg)
	default:
		return errors.Errorf("unknown config format '%s'", format)
	}

	if err != nil {
		return errors.Wrapf(err, "formatting config (%s)", format)
	}

	_, err = fmt.Fprintf(w, "\n%s\n", data)
	return err
}

func defaultConfig() *Config {
	return &Config{
		Azure: AzureConfig{
			CloudName: defaultCloudName,
			UserAgent: release.NAME,
			Interval:  defaultInterval,
		},
	}
}
