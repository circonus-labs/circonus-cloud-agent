// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package gcpservice

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/circonus-labs/circonus-cloud-agent/internal/services/gcpservice/collectors"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Instance defines a specific gcp service instance for collecting metrics
type Instance struct {
	cfg        *Config
	ctx        context.Context
	logger     zerolog.Logger
	check      *circonus.Check
	collectors []collectors.Collector
	lastStart  *time.Time
	baseTags   circonus.Tags
	running    bool
	sync.Mutex
}

// Start runs the instance based on the configured interval
func (inst *Instance) Start() error {
	interval := time.Duration(inst.cfg.GCP.Interval) * time.Minute
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-inst.ctx.Done():
			return nil
		case <-ticker.C:
			inst.Lock()
			if inst.lastStart != nil && time.Since(*inst.lastStart) < interval {
				inst.Unlock()
				continue
			}
			if inst.running {
				inst.Unlock()
				inst.logger.Warn().Msg("collection in progress, not starting another")
				continue
			}

			// calculate one timeseries range for all requests from collectors
			start := time.Now()
			var delta time.Duration
			if inst.lastStart == nil {
				delta = interval * 2
			} else {
				delta = start.Sub(*inst.lastStart) + 2*time.Minute
			}
			tsEnd := start
			tsStart := tsEnd.Add(-delta)
			inst.logger.Info().Time("start", tsStart).Time("end", tsEnd).Str("delta", delta.String()).Msg("collection timeseries range")

			inst.lastStart = &start
			inst.running = true
			inst.Unlock()

			go func() {
				for _, c := range inst.collectors {
					if err := c.Collect(tsStart, tsEnd, inst.cfg.GCP.projectID, inst.cfg.GCP.credentialData, inst.baseTags); err != nil {
						inst.check.ReportError(errors.WithMessage(err, fmt.Sprintf("id: %s, collector: %s", inst.cfg.ID, c.ID())))
						inst.logger.Warn().Err(err).Str("collector", c.ID()).Msg("collecting telemetry")
						// need to determine which errors from the various
						// cloud service providers are fatal vs retry vs
						// wait for next iteration
					}
					if inst.done() {
						break
					}
				}
				inst.Lock()
				inst.running = false
				inst.Unlock()
				inst.logger.Debug().Str("duration", time.Since(start).String()).Msg("collection complete")
			}()
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
