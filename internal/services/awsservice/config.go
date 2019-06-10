// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package awsservice

// AWS service configuration

import (
	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/circonus-labs/circonus-cloud-agent/internal/services/awsservice/collectors"
)

// Config defines an AWS service instance configuration
// NOTE: warning - ID must be thought of as immutable - if it changes a new check will be created
type Config struct {
	ID       string                 `json:"id" toml:"id" yaml:"id"`                   // unique id for this service client instance, no spaces (ties several things together, short and immutable - logging, check search/create, tags, etc.)
	Regions  []AWSRegion            `json:"regions" toml:"regions" yaml:"regions"`    // list of region specific configurations
	AWS      AWS                    `json:"aws" toml:"aws" yaml:"aws"`                // REQUIRED, aws credentials
	Circonus circonus.ServiceConfig `json:"circonus" toml:"circonus" yaml:"circonus"` // REQUIRED, circonus config: api credentials, check, broker, etc.
	Period   string                 `json:"period" toml:"period" yaml:"period"`       // 'basic' or 'detailed'
	Tags     circonus.Tags          `json:"tags" toml:"tags" yaml:"tags"`             // global tags, added to all metrics
}

// AWSRegion defines a specific aws region from which to collect metrics
type AWSRegion struct {
	Name     string                    `json:"name" toml:"name" yaml:"name"`             // e.g. us-east-1
	Services []collectors.AWSCollector `json:"services" toml:"services" yaml:"services"` // which services to collectc metrics for in this region
	Tags     circonus.Tags             `json:"tags" toml:"tags" yaml:"tags"`             // region tags (default region:Name)
}

// AWS defines the credentials to use for AWS
type AWS struct {
	//
	// Runs in ONE of two ways: shared or local
	//   shared - multiple different sets of credentials - use AccessKeyID and SecretAccessKey from an IAM role
	//   local - use a role in a local credentials file
	//
	// Mode is swithed based on which is configured, configure only ONE pair of authentication settings
	//
	// see: https://docs.aws.amazon.com/sdk-for-go/api/aws/credentials/#Value
	// to be used with: https://docs.aws.amazon.com/sdk-for-go/api/aws/credentials/#NewStaticCredentials
	AccessKeyID     string `json:"access_key_id" toml:"access_key_id" yaml:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key" toml:"secret_access_key" yaml:"secret_access_key"`
	// SessionToken    string `json:"session_token" toml:"session_token" yaml:"session_token"`
	// ProviderName    string `json:"provider_name" toml:"provider_name" yaml:"provider_name"`
	//
	// Local mode - non-shared version, running a local instance of circonus-cloud-agent
	// use with https://docs.aws.amazon.com/sdk-for-go/api/aws/credentials/#NewSharedCredentials
	Role            string `json:"role" toml:"role" yaml:"role"`
	CredentialsFile string `json:"credentials_file" toml:"credentials_file" yaml:"credentials_file"`
}
