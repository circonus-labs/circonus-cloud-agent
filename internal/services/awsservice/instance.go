// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package awsservice

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/circonus-labs/circonus-cloud-agent/internal/config"
	"github.com/circonus-labs/circonus-cloud-agent/internal/release"
	"github.com/circonus-labs/circonus-cloud-agent/internal/services/awsservice/collectors"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Instance AWS SDK/API Instance for fetching cloudwatch metrics and forwarding them to Circonus
// Note: a Instance has a 1:1 relation with aws:circ - each Instance has (or, may have)
// a different set of aws and/or circonus credentials.
type Instance struct {
	cfg        *Config
	regionCfg  *AWSRegion
	ctx        context.Context
	interval   uint
	logger     zerolog.Logger
	check      *circonus.Check
	period     int64
	collectors []collectors.Collector
	baseTags   circonus.Tags
	running    bool
	sync.Mutex
}

// initInstances creates a new set of region Instances for each configuration
// file (.json,.toml,.yaml) found in the passed confDir
func (svc *AWSService) initInstances(confDir string) error {
	if confDir == "" {
		return errors.New("invalid config dir (empty)")
	}

	entries, err := ioutil.ReadDir(confDir)
	if err != nil {
		return errors.Wrap(err, "reading AWS config dir")
	}

	for _, entry := range entries {
		entry := entry
		if entry.IsDir() {
			continue
		}
		if !strings.Contains(".json|.toml|.yaml", filepath.Ext(entry.Name())) {
			continue
		}

		var cfg Config
		if err := config.LoadConfigFile(path.Join(confDir, entry.Name()), &cfg); err != nil {
			svc.logger.Error().Err(err).Str("file", entry.Name()).Msg("loading config, skipping")
			continue
		}

		if cfg.ID == "" {
			svc.logger.Error().Str("file", entry.Name()).Msg("invalid config ID (empty), skipping")
			continue
		}
		if strings.Contains(cfg.ID, " ") {
			svc.logger.Error().Str("file", entry.Name()).Msg("invalid config ID (contains spaces), skipping")
			continue
		}
		if len(cfg.Regions) == 0 {
			svc.logger.Error().Str("file", entry.Name()).Msg("invalid config regions (empty), skipping")
		}

		// based on cfg.Period - collect every 1min for 'detailed' or every 5min for 'basic'
		period := 300
		if cfg.Period == "detailed" {
			period = 60
		}
		// used to control how many samples we request - calclulating start from
		// time.Now (e.g. time.Now().Add(- (interval * time.Second))). desired
		// number of samples is three. if exactly three * period is used,
		// cloudwatch sdk will often respond with only the last two samples.
		// so use 3 * period, plus a little extra cushion.
		interval := (period * 3) + (period / 2)

		for _, regionConfig := range cfg.Regions {
			regionConfig := regionConfig
			instance := &Instance{
				cfg:       &cfg,
				regionCfg: &regionConfig,
				ctx:       svc.groupCtx,
				interval:  uint(interval),
				logger:    svc.logger.With().Str("id", cfg.ID).Str("region", regionConfig.Name).Logger(),
				period:    int64(period),
			}
			instance.logger.Debug().Str("aws_region", regionConfig.Name).Msg("initialized client instance for region")

			checkConfig := &circonus.Config{
				ID:            fmt.Sprintf("aws_%s_%s", cfg.ID, regionConfig.Name),
				DisplayName:   fmt.Sprintf("aws %s %s /%s", cfg.ID, regionConfig.Name, release.NAME),
				CheckBundleID: cfg.Circonus.CID,
				APIKey:        cfg.Circonus.Key,
				APIApp:        cfg.Circonus.App,
				APIURL:        cfg.Circonus.URL,
				Debug:         cfg.Circonus.Debug,
				Logger:        instance.logger,
				Tags:          fmt.Sprintf("%s:aws,aws_region:%s", release.NAME, regionConfig.Name),
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
				instance.logger.Error().Err(err).Msg("creating Circonus Check instance, skipping")
				continue
			}
			instance.check = chk

			ms, err := collectors.New(instance.ctx, instance.check, regionConfig.Services, instance.logger)
			if err != nil {
				instance.logger.Warn().Err(err).Msg("setting up aws metric services")
				continue
			}
			instance.collectors = ms

			svc.instances = append(svc.instances, instance)
		}
	}

	if len(svc.instances) == 0 {
		return errors.New("no valid AWS configs found")
	}

	return nil
}

// Start metric collections based on the configured interval - intended to be run in a goroutine (e.g. errgroup)
func (inst *Instance) Start() error {
	ticker := time.NewTicker(time.Duration(inst.period) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-inst.ctx.Done():
			return nil
		case <-ticker.C:
			inst.logger.Debug().Msg("metric collection triggered")
			sess, err := inst.createSession(inst.regionCfg.Name)
			if err != nil {
				inst.logger.Warn().Err(err).Msg("creating AWS SDK session")
				continue
			}

			inst.Lock()
			if inst.running {
				inst.Unlock()
				inst.logger.Warn().Msg("collection already in progress, not starting another")
				continue
			}
			inst.running = true
			inst.Unlock()

			end := time.Now()
			start := end.Add(-(time.Duration(inst.interval) * time.Second))
			timespan := collectors.MetricTimespan{
				Start:  start,
				End:    end,
				Period: inst.period,
			}
			for _, c := range inst.collectors {
				if err := c.Collect(sess, timespan, inst.baseTags); err != nil {
					inst.check.ReportError(errors.WithMessage(err, fmt.Sprintf("id: %s, collector: %s", inst.cfg.ID, c.ID())))
					inst.logger.Warn().Err(err).Str("collector", c.ID()).Msg("collecting telemetry")
				}
			}

			inst.Lock()
			inst.running = false
			inst.Unlock()
		}
	}
}

// done is a utility routine to check the context, returns true if done
func (inst *Instance) done() bool {
	select {
	case <-inst.ctx.Done():
		inst.logger.Debug().Msg("context done, exiting")
		return true
	default:
		return false
	}
}

// createSession returns a new aws session using configured aws information
func (inst *Instance) createSession(region string) (*session.Session, error) {
	var creds *credentials.Credentials

	switch {
	case inst.cfg.AWS.Role != "":
		creds = credentials.NewSharedCredentials(
			inst.cfg.AWS.CredentialsFile,
			inst.cfg.AWS.Role)
	case inst.cfg.AWS.AccessKeyID != "":
		creds = credentials.NewStaticCredentials(
			inst.cfg.AWS.AccessKeyID,
			inst.cfg.AWS.SecretAccessKey,
			"")
	default:
		return nil, errors.New("invalid AWS credentils configuration")
	}

	cfg := &aws.Config{Credentials: creds}
	if region != "" && region != "global" {
		cfg.Region = aws.String(region)
	}

	return session.NewSession(cfg)
}
