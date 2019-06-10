// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package azureservice

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/circonus-labs/circonus-cloud-agent/internal/config"
	"github.com/circonus-labs/circonus-cloud-agent/internal/config/defaults"
	"github.com/circonus-labs/circonus-cloud-agent/internal/release"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

const (
	// KeyEnabled toggles whether the Azure module is active or not
	KeyEnabled = "azure.enabled"
	// DefaultEnabled defines the default setting
	DefaultEnabled = false

	// KeyConfDir defines the azure configuration directory
	KeyConfDir = "azure.conf_dir"

	// KeyConfExample shows an example configuration
	KeyConfExample = "azure.config_example"
)

var (
	// DefaultConfDir is the default location of azure configuration files
	DefaultConfDir = path.Join(defaults.EtcPath, "azure.d")
)

// AzureService defines the Azure cloud service client
type AzureService struct {
	enabled   bool
	group     *errgroup.Group
	groupCtx  context.Context
	instances []*Instance
	logger    zerolog.Logger
}

// New returns an Azure cloud service metric collection client
func New(ctx context.Context) (*AzureService, error) {
	if fmt := viper.GetString(KeyConfExample); fmt != "" {
		if err := showExampleConfig(fmt, os.Stdout); err != nil {
			log.Error().Err(err).Msg("generating example configuration")
		}
		os.Exit(0)
	}

	g, gctx := errgroup.WithContext(ctx)
	svc := AzureService{
		enabled:   viper.GetBool(KeyEnabled),
		group:     g,
		groupCtx:  gctx,
		instances: make([]*Instance, 0),
		logger:    log.With().Str("pkg", "azure").Logger(),
	}

	if !svc.enabled {
		return &svc, nil
	}

	svc.logger.Debug().Msg("Azure enabled, initializing client")

	confDir := viper.GetString(KeyConfDir)
	if confDir == "" {
		confDir = DefaultConfDir
	}

	if err := svc.initInstances(confDir); err != nil {
		return nil, errors.Wrap(err, "initializing Azure metric collector instances(s)")
	}

	if len(svc.instances) == 0 {
		svc.logger.Info().Msg("disabling Azure, no metric collectors initialized")
		svc.enabled = false
	}

	return &svc, nil
}

// Enabled indicates whether the Azure service is enabled
func (svc *AzureService) Enabled() bool {
	return svc.enabled
}

// Scan checks the service config directory for configurations and loads them
func (svc *AzureService) Scan() error {
	return errors.New("not implemented")
}

// Start begins collecting metrics from Azure
func (svc *AzureService) Start() error {
	if !svc.enabled {
		svc.logger.Info().Msg("Azure client disabled, not starting")
		return nil
	}

	svc.logger.Info().Msg("Azure client starting")

	for _, instance := range svc.instances {
		inst := instance
		svc.group.Go(inst.Start)
	}

	go func() {
		if err := svc.group.Wait(); err != nil {
			svc.logger.Error().Err(err).Msg("waiting for service group")
		}
	}()

	return svc.group.Wait()
}

func (svc *AzureService) initInstances(confDir string) error {
	if confDir == "" {
		return errors.New("invalid config dir (empty)")
	}

	entries, err := ioutil.ReadDir(confDir)
	if err != nil {
		return errors.Wrap(err, "reading AWS config dir")
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.Contains(".json|.toml|.yaml", filepath.Ext(entry.Name())) {
			continue
		}

		cfgFile := path.Join(confDir, entry.Name())
		instance, err := svc.instanceFromConfig(cfgFile)
		if err != nil {
			svc.logger.Error().Err(err).Str("config_file", cfgFile).Msg("skipping")
			continue
		}

		svc.instances = append(svc.instances, instance)
	}

	if len(svc.instances) == 0 {
		return errors.New("no valid Azure configs found")
	}

	return nil
}

func (svc *AzureService) instanceFromConfig(cfgFile string) (*Instance, error) {
	var cfg Config
	if err := config.LoadConfigFile(cfgFile, &cfg); err != nil {
		return nil, errors.Wrap(err, "loading config file")
	}

	if cfg.ID == "" {
		return nil, errors.New("invalid config ID (empty)")
	}
	if strings.Contains(cfg.ID, " ") {
		return nil, errors.New("invalid config ID (contains spaces)")
	}

	instance := &Instance{
		cfg:    &cfg,
		ctx:    svc.groupCtx,
		logger: svc.logger.With().Str("id", cfg.ID).Logger(),
	}

	// handle config settings and/or overlay defaults
	if cfg.Azure.Interval < defaultInterval {
		cfg.Azure.Interval = defaultInterval
	}
	if cfg.Azure.CloudName == "" {
		cfg.Azure.CloudName = defaultCloudName
	}
	if cfg.Azure.UserAgent == "" {
		cfg.Azure.UserAgent = release.NAME
	}

	sm, err := instance.getSubscriptionMeta()
	if err != nil {
		return nil, errors.Wrap(err, "init instance, azure subscription info")
	}

	instance.logger.Info().Str("subscription", sm.Name).Msg("creating instance")

	checkConfig := &circonus.Config{
		ID:            fmt.Sprintf("azure_%s", cfg.ID),
		DisplayName:   fmt.Sprintf("azure %s %s /%s", cfg.ID, sm.Name, release.NAME),
		CheckBundleID: cfg.Circonus.CID,
		APIKey:        cfg.Circonus.Key,
		APIApp:        cfg.Circonus.App,
		APIURL:        cfg.Circonus.URL,
		Debug:         cfg.Circonus.Debug,
		Logger:        instance.logger,
		Tags:          fmt.Sprintf("%s:azure", release.NAME),
	}
	if len(cfg.Tags) > 0 { // if top-level tags are configured, add them to check
		tags := make([]string, len(cfg.Tags))
		for idx, tag := range cfg.Tags {
			tags[idx] = tag.Category + ":" + tag.Value
		}
		checkConfig.Tags += "," + strings.Join(tags, ",")
	}

	chk, err := circonus.NewCheck(checkConfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating/retrieving circonus check")
	}
	instance.check = chk

	return instance, nil
}
