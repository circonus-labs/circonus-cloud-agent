// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package gcpservice

import (
	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/circonus-labs/circonus-cloud-agent/internal/services/gcpservice/collectors"
)

// Config defines the options for a gcp service instance.
type Config struct {
	ID       string                 `json:"id" toml:"id" yaml:"id"`                   // unique id for this service client instance, no spaces (ties several things together, short and immutable - logging, check search/create, tags, etc.)
	Circonus circonus.ServiceConfig `json:"circonus" toml:"circonus" yaml:"circonus"` // REQUIRED circonus config: api credentials, check, broker, etc.
	Tags     circonus.Tags          `json:"tags" toml:"tags" yaml:"tags"`             // global tags, added to all metrics
	GCP      GCPConfig              `json:"gcp" toml:"gcp" yaml:"gcp"`                // REQUIRED gcp configuration
}

// GCPConfig holds the gcp specific configuration options.
type GCPConfig struct {
	Collectors      []collectors.GCPCollector `json:"services" toml:"services" yaml:"services"`                         // which services to collectc metrics for in this region
	CredentialsFile string                    `json:"credentials_file" toml:"credentials_file" yaml:"credentials_file"` // REQUIRED gcp service account credentials file (if it does not begin with filepath.Separator, the default etc path - relative to where circonus-cloud-agentd is running - will be used, allowing relative paths)
	projectID       string                    // extracted from CredentialsFile
	projectName     string                    // from project meta data retrieved using gcp api, used for check display name
	credentialData  []byte                    // loaded from CredentialsFile
	Interval        int                       `json:"collect_interval" toml:"collect_interval" yaml:"collect_interval"` // DEFAULT: 5 - How often to collect metrics.(minutes >= 5)
}

var (
	defaultInterval = 5
)
