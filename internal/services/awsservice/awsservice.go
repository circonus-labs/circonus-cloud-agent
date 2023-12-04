// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package awsservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/circonus-labs/circonus-cloud-agent/internal/config/defaults"
	"github.com/circonus-labs/circonus-cloud-agent/internal/services/awsservice/collectors"
	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	yaml "gopkg.in/yaml.v2"
)

const (
	// KeyEnabled toggles whether the AWS module is active or not.
	KeyEnabled = "aws.enabled"
	// DefaultEnabled defines the default setting.
	DefaultEnabled = false

	// KeyConfDir defines the aws configuration directory.
	KeyConfDir = "aws.conf_dir"

	// KeyConfExample shows an example configuration.
	KeyConfExample = "aws.config_example"
)

var (
	// DefaultConfDir is the default location of aws configuration files.
	DefaultConfDir = path.Join(defaults.EtcPath, "aws.d")
)

// AWSService defines the AWS cloud service client.
type AWSService struct {
	groupCtx  context.Context
	group     *errgroup.Group
	instances []*Instance
	logger    zerolog.Logger
	enabled   bool
}

// New returns an AWS cloud service metric collector.
func New(ctx context.Context) (*AWSService, error) {
	if fmt := viper.GetString(KeyConfExample); fmt != "" {
		if err := showExampleConfig(fmt, os.Stdout); err != nil {
			log.Error().Err(err).Msg("generating example configuration")
		}
		os.Exit(0)
	}

	g, gctx := errgroup.WithContext(ctx)
	svc := AWSService{
		enabled:   viper.GetBool(KeyEnabled),
		group:     g,
		groupCtx:  gctx,
		instances: make([]*Instance, 0),
		logger:    log.With().Str("pkg", "aws").Logger(),
	}

	if !svc.enabled {
		return &svc, nil
	}

	svc.logger.Debug().Msg("AWS enabled, initializing client")

	confDir := viper.GetString(KeyConfDir)
	if confDir == "" {
		confDir = DefaultConfDir
	}

	if err := svc.initInstances(confDir); err != nil {
		return nil, errors.Wrap(err, "initializing AWS metric collector instances(s)")
	}

	if len(svc.instances) == 0 {
		svc.logger.Info().Msg("disabling AWS, no metric collectors initialized")
		svc.enabled = false
	}

	return &svc, nil
}

// Enabled indicates whether the AWS service is enabled.
func (svc *AWSService) Enabled() bool {
	return svc.enabled
}

// Scan checks the service config directory for configurations and loads them.
func (svc *AWSService) Scan() error {
	return errors.New("not implemented")
}

// Start begins collecting metrics from AWS service.
func (svc *AWSService) Start() error {
	if !svc.enabled {
		svc.logger.Info().Msg("AWS client disabled, not starting")
		return nil
	}
	svc.logger.Info().Msg("AWS client starting")

	// start the aws service instance(s)
	for _, instance := range svc.instances {
		inst := instance
		svc.group.Go(inst.Start)
	}

	go func() {
		if err := svc.group.Wait(); err != nil {
			svc.logger.Warn().Err(err).Msg("waiting for service group")
		}
	}()

	return svc.group.Wait()
}

func showExampleConfig(format string, w io.Writer) error {
	var err error
	var data []byte

	cfg := &Config{}

	cc, err := collectors.ConfigExample()
	if err != nil {
		return errors.Wrap(err, "generating example service collectors config")
	}
	cfg.Regions = []AWSRegion{
		{
			Name:     "us-east-1",
			Services: cc,
		},
	}

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
