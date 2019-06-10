// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package azureservice

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Instance Azure SDK/API Instance for fetching metrics and forwarding them to Circonus
// Note: a Instance has a 1:1 relation with azure:circ - each Instance has (or, may have)
// a different set of azure and/or circonus credentials.
type Instance struct {
	cfg      *Config
	ctx      context.Context
	logger   zerolog.Logger
	check    *circonus.Check
	baseTags circonus.Tags
	running  bool
	sync.Mutex
}

// Start runs the instance based on the configured interval
func (inst *Instance) Start() error {
	ticker := time.NewTicker(time.Duration(inst.cfg.Azure.Interval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-inst.ctx.Done():
			return nil
		case <-ticker.C:
			inst.logger.Debug().Msg("metric collection triggered")
			inst.Lock()
			if inst.running {
				inst.Unlock()
				inst.logger.Warn().Msg("collection already in progress, not starting another")
				continue
			}
			inst.running = true
			inst.Unlock()

			endTime := time.Now().UTC()

			err := inst.collect(endTime)
			if err != nil {
				inst.check.ReportError(errors.WithMessage(err, fmt.Sprintf("id: %s", inst.cfg.ID)))
				inst.logger.Warn().Err(err).Msg("collecting metrics")
				// if fatal return the error
				// need to determine which errors from the various cloud service providers are fatal vs retry vs wait for next iteration
			}

			// if viper.GetBool(config.KeyPipeSubmits) {
			// 	err = inst.collectWithPipe(endTime)
			// } else {
			//  err = inst.collect(endTime)
			// }

			inst.Lock()
			inst.running = false
			inst.Unlock()
		}
	}
}

// collect metrics from Azure and forward to Circonus using buffer
func (inst *Instance) collect(endTime time.Time) error {
	// NOTE: this model needs to be used, so submission requests will have
	// a Content-Length while streaming JSON data:
	//
	// 1. Create a buffer
	// 2. For each resource
	//    a. Collect resource metrics (write into buffer)
	//    b. Submit metrics (read from buffer)
	//    c. Reset buffer (so it can be re-used for next resource)
	//
	// Given there is no way to know how many metrics/samples will
	// accutally be received from any given resource. Safer to collect
	// from each resource and submit immediately/independently.
	// Rather than, collecting all metrics/samples into one buffer and
	// submitting as one potentially very large, fragile batch.

	auth, err := inst.authorize()
	if err != nil {
		return errors.Wrap(err, "authorize, subscription meta")
	}

	resources, err := inst.getResources(auth)
	if err != nil {
		return errors.Wrap(err, "resource list")
	}

	if inst.done() {
		return nil
	}

	var buf bytes.Buffer
	buf.Grow(32678)

	for _, resource := range resources {
		if inst.done() {
			break
		}

		err := inst.getResourceMetrics(&buf, auth, resource.ID, endTime, resource.Tags)
		if err != nil {
			inst.check.ReportError(errors.WithMessage(err, fmt.Sprintf("id: %s, resource_id: %s", inst.cfg.ID, resource.ID)))
			inst.logger.Warn().Err(err).Str("resource_id", resource.ID).Msg("collecting metrics")
			// NOTE: do not 'continue' here, fall-through so that any metric
			// samples collected prior to the error get submitted.
			//

		}

		if buf.Len() == 0 {
			inst.logger.Warn().Str("resource_id", resource.ID).Msg("no telemetry to submit")
			continue
		}

		inst.logger.Debug().Str("resource_id", resource.ID).Msg("submitting telemetry")
		if err := inst.check.SubmitMetrics(&buf); err != nil {
			inst.check.ReportError(errors.WithMessage(err, fmt.Sprintf("id: %s, resource_id: %s", inst.cfg.ID, resource.ID)))
			inst.logger.Error().Err(err).Str("resource_id", resource.ID).Msg("submitting telemetry")
		}

		buf.Reset()
	}

	// TODO: submit run stats (e.g. buf.Reset(); write run metrics, submit run metrics)

	return nil
}

// collectWithPipe metrics from Azure and forward to Circonus using pipes
func (inst *Instance) collectWithPipe(endTime time.Time) error {
	// NOTE: Unable to use pipes at the moment. ATS &| broker cannot
	// handle PUT|POST requests without a Content-Length header which is,
	// of course, not possible with a pipe...

	auth, err := inst.authorize()
	if err != nil {
		return errors.Wrap(err, "authorize, subscription meta")
	}

	resources, err := inst.getResources(auth)
	if err != nil {
		return errors.Wrap(err, "resource list")
	}

	if inst.done() {
		return nil
	}

	var wg sync.WaitGroup
	pr, pw := io.Pipe()
	wg.Add(1)

	go func() {
		defer wg.Done()

		for _, resource := range resources {
			err := inst.getResourceMetrics(pw, auth, resource.ID, endTime, resource.Tags)
			if err != nil {
				inst.logger.Warn().Err(err).Str("resource_id", resource.ID).Msg("collecting metrics")
			}

			if inst.done() {
				break
			}
		}

		// TODO: submit run stats, write run metrics to pw

		if err := pw.Close(); err != nil {
			inst.logger.Error().Err(err).Msg("closing pipe writer")
		}
	}()

	inst.logger.Debug().Msg("starting metric submission")
	if err := inst.check.SubmitMetrics(pr); err != nil {
		inst.logger.Error().Err(err).Msg("submitting metrics")
	}

	wg.Wait()

	// queue up errors
	// submit each error here in separate request

	return nil
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
