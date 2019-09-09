// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package gcpservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/circonus-labs/circonus-cloud-agent/internal/config"
	"github.com/circonus-labs/circonus-cloud-agent/internal/config/defaults"
	"github.com/circonus-labs/circonus-cloud-agent/internal/release"
	"github.com/circonus-labs/circonus-cloud-agent/internal/services/gcpservice/collectors"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

const (
	// KeyEnabled toggles whether the GCP module is active or not
	KeyEnabled = "gcp.enabled"
	// DefaultEnabled defines the default setting
	DefaultEnabled = false

	// KeyConfDir defines the gcp configuration directory
	KeyConfDir = "gcp.conf_dir"

	// KeyConfExample shows an example configuration
	KeyConfExample = "gcp.config_example"
)

var (
	// DefaultConfDir is the default location of gcp configuration files
	DefaultConfDir = path.Join(defaults.EtcPath, "gcp.d")
)

// GCPService defines the GCP cloud service client
type GCPService struct {
	enabled   bool
	group     *errgroup.Group
	groupCtx  context.Context
	instances []*Instance
	logger    zerolog.Logger
}

// New returns a GCP metric collector service
func New(ctx context.Context) (*GCPService, error) {
	if fmt := viper.GetString(KeyConfExample); fmt != "" {
		if err := showExampleConfig(fmt, os.Stdout); err != nil {
			log.Error().Err(err).Msg("generating example configuration")
		}
		os.Exit(0)
	}

	g, gctx := errgroup.WithContext(ctx)
	svc := GCPService{
		enabled:   viper.GetBool(KeyEnabled),
		group:     g,
		groupCtx:  gctx,
		instances: make([]*Instance, 0),
		logger:    log.With().Str("pkg", "gcp").Logger(),
	}

	if !svc.enabled {
		return &svc, nil
	}

	svc.logger.Debug().Msg("enabled, initializing client")

	confDir := viper.GetString(KeyConfDir)
	if confDir == "" {
		confDir = DefaultConfDir
	}

	if err := svc.initInstances(confDir); err != nil {
		return nil, errors.Wrap(err, "initializing telemetry collector(s)")
	}

	if len(svc.instances) == 0 {
		svc.logger.Info().Msg("disabling, no telemetry collector(s) initialized")
		svc.enabled = false
	}

	return &svc, nil
}

// Enabled indicates whether the GCP service is enabled
func (svc *GCPService) Enabled() bool {
	return svc.enabled
}

// Scan checks the service config directory for configurations and loads them
func (svc *GCPService) Scan() error {
	if !svc.enabled {
		svc.logger.Info().Msg("client disabled, not checking for configurations")
		return nil
	}
	svc.logger.Info().Msg("client checking for configuration(s)")
	return errors.New("not implemented")
}

// Start begins collecting metrics from GCP
func (svc *GCPService) Start() error {
	if !svc.enabled {
		svc.logger.Info().Msg("client disabled, not starting")
		return nil
	}
	svc.logger.Info().Msg("client starting")

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

func (svc *GCPService) initInstances(confDir string) error {
	if confDir == "" {
		return errors.New("invalid config dir (empty)")
	}

	entries, err := ioutil.ReadDir(confDir)
	if err != nil {
		return errors.Wrap(err, "reading GCP config dir")
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
		return errors.New("no valid gcp configs found")
	}

	return nil
}

func (svc *GCPService) instanceFromConfig(cfgFile string) (*Instance, error) {
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

	if cfg.GCP.CredentialsFile == "" {
		return nil, errors.New("invalid GCP credentials file (empty)")
	}

	cfgFile, err := config.VerifyFile(instance.cfg.GCP.CredentialsFile)
	if err != nil {
		return nil, errors.Wrap(err, "invalid GCP credentials file")
	}
	if !strings.HasPrefix(cfgFile, string(filepath.Separator)) {
		cfgFile = filepath.Join(defaults.EtcPath, cfgFile)
	}
	instance.cfg.GCP.CredentialsFile = cfgFile

	data, err := ioutil.ReadFile(instance.cfg.GCP.CredentialsFile)
	if err != nil {
		return nil, errors.Wrap(err, "loading GCP credentials file")
	}
	instance.cfg.GCP.credentialData = data

	var v struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	instance.cfg.GCP.projectID = v.ProjectID

	if err := instance.loadProjectMeta(); err != nil {
		return nil, err
	}

	if cfg.GCP.Interval < defaultInterval {
		cfg.GCP.Interval = defaultInterval
	}

	checkConfig := &circonus.Config{
		ID:            "gcp_" + instance.cfg.ID,
		DisplayName:   fmt.Sprintf("%s %s %s/gcp", instance.cfg.ID, instance.cfg.GCP.projectName, release.NAME),
		CheckBundleID: cfg.Circonus.CID,
		APIKey:        cfg.Circonus.Key,
		APIApp:        cfg.Circonus.App,
		APIURL:        cfg.Circonus.URL,
		Debug:         cfg.Circonus.Debug,
		TraceMetrics:  cfg.Circonus.TraceMetrics,
		Logger:        instance.logger,
		Tags:          release.NAME + ":gcp",
	}
	if len(instance.cfg.Tags) > 0 { // if top-level tags are configured, add them to check
		tags := make([]string, len(instance.cfg.Tags))
		for idx, tag := range instance.cfg.Tags {
			tags[idx] = tag.Category + ":" + tag.Value
		}
		checkConfig.Tags += "," + strings.Join(tags, ",")
	}
	if len(instance.baseTags) > 0 { // if project-level tags were defined, add them to check
		tags := make([]string, len(instance.baseTags))
		for idx, tag := range instance.baseTags {
			tags[idx] = tag.Category + ":" + tag.Value
		}
		checkConfig.Tags += "," + strings.Join(tags, ",")
	}

	chk, err := circonus.NewCheck("gcp", checkConfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating/retrieving circonus check")
	}
	instance.check = chk

	// initialize collectors
	pollingInterval := time.Duration(cfg.GCP.Interval) * time.Minute
	ms, err := collectors.New(instance.ctx, instance.check, cfg.GCP.Collectors, pollingInterval, instance.logger)
	if err != nil {
		return nil, err
	}
	instance.collectors = ms

	return instance, nil
}

func showExampleConfig(format string, w io.Writer) error {
	var err error
	var data []byte

	cfg := &Config{}
	cfg.GCP.Interval = 5

	cc, err := collectors.ConfigExample()
	if err != nil {
		return errors.Wrap(err, "generating example service collectors config")
	}
	cfg.GCP.Collectors = cc

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
