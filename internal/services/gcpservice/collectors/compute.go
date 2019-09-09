// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

// Compute holds definition for the gcp gce collector
type Compute struct {
	common
}

func newCompute(ctx context.Context, check *circonus.Check, cfg *GCPCollector, interval time.Duration, logger zerolog.Logger) (Collector, error) {
	if ctx == nil {
		return nil, errors.New("invalid context (nil)")
	}
	if check == nil {
		return nil, errors.New("invalid check (nil)")
	}
	if cfg == nil {
		return nil, errors.New("invalid config (nil)")
	}
	c := &Compute{
		common: newCommon(ctx, check, cfg, interval, logger),
	}
	c.logger.Debug().Msg("initialized")
	return c, nil
}

// Collect telemetry from gce instance resources in project
func (c *Compute) Collect(timeseriesStart, timeseriesEnd time.Time, projectID string, creds []byte, baseTags circonus.Tags) error {

	if !c.enabled {
		if c.disableCause == "" || c.disableTime == nil {
			return nil
		}
		if time.Since(*c.disableTime) >= 1*time.Hour {
			c.enabled = true // RE-enable to try again, since it is configured to be enabled
		} else {
			c.logger.Warn().
				Time("disable_time", *c.disableTime).
				Str("disable_cause", c.disableCause).
				Msg("collector has been disabled due to error")
			return nil
		}
	}

	runStart := timeseriesEnd
	c.tsEnd = timeseriesEnd
	c.tsStart = timeseriesStart

	instanceList, err := c.getInstanceList(projectID, creds)
	if err != nil {
		return err
	}

	if len(instanceList) == 0 {
		c.logger.Debug().Msg("zero instances to process")
		return nil
	}

	var buf bytes.Buffer
	buf.Grow(32768)
	c.logger.Debug().Int("instances", len(instanceList)).Msg("processing instances")
	for _, info := range instanceList {
		instStart := time.Now()
		instLogger := c.logger.With().Str("region", info.region).Str("instance", info.name).Logger()

		metricFilter := fmt.Sprintf(`metric.labels.instance_name = "%s"`, info.name)
		if err := c.processMetrics(projectID, metricFilter, creds, &buf, baseTags); err != nil {
			instLogger.Warn().Err(err).Msg("collecting instance metrics")
		}
		instLogger.Info().Str("duration", time.Since(instStart).String()).Msg("instance collect end")

		if buf.Len() == 0 {
			instLogger.Warn().Msg("no telemetry to submit")
			continue
		}

		submitStart := time.Now()
		if err := c.check.SubmitMetrics(&buf); err != nil {
			c.check.ReportError(errors.WithMessage(err, fmt.Sprintf("collector: %s", c.ID())))
			instLogger.Error().Err(err).Msg("submitting telemetry")
		}
		instLogger.Info().Str("duration", time.Since(submitStart).String()).Msg("instance submit end")

		buf.Reset()
		instLogger.Info().Str("duration", time.Since(instStart).String()).Msg("instance run end")

		if c.done() {
			break
		}
	}

	c.logger.Info().Str("duration", time.Since(runStart).String()).Msg("compute run end")

	return nil
}

// instanceInfo holds the meta data per gce instance used to collect telemetry metrics
type instanceInfo struct {
	region string
	name   string
}

// getInstanceList returns a list of gce instances
func (c *Compute) getInstanceList(projectID string, creds []byte) ([]instanceInfo, error) {
	emptyList := []instanceInfo{}

	filter := ""
	if c.filter.Expression != "" {
		filter = c.filter.Expression
	}
	if filter == "" && len(c.filter.Labels) > 0 {
		var expressions []string
		for k, v := range c.filter.Labels {
			if k == "" || v == "" {
				continue
			}
			expressions = append(expressions, fmt.Sprintf(`(labels.%s = "%s")`, k, v))
		}
		if len(expressions) > 0 {
			filter = strings.Join(expressions, "")
		}
	}

	c.logger.Debug().Str("filter", filter).Msg("gce instance filter")

	// https://godoc.org/google.golang.org/api/compute/v1#InstancesService
	conf, err := google.JWTConfigFromJSON(creds, "https://www.googleapis.com/auth/compute.readonly")
	if err != nil {
		return emptyList, errors.Wrap(err, "token for compute")
	}
	client := conf.Client(c.ctx)

	// svc, err := compute.New(client)
	svc, err := compute.NewService(c.ctx, option.WithHTTPClient(client))
	if err != nil {
		return emptyList, errors.Wrap(err, "compute client")
	}

	instanceList := make([]instanceInfo, 0)

	// https://godoc.org/google.golang.org/api/compute/v1#NewInstancesService
	isvc := compute.NewInstancesService(svc)
	nextPageToken := ""
	for {
		// https://godoc.org/google.golang.org/api/compute/v1#InstancesService.AggregatedList
		iac := isvc.AggregatedList(projectID)
		ial, err := iac.Context(c.ctx).PageToken(nextPageToken).Filter(filter).Do()
		if err != nil {
			return instanceList, errors.Wrap(err, "instances aggregated list")
		}

		for region, item := range ial.Items {
			if len(item.Instances) == 0 {
				continue
			}
			for _, gceInstance := range item.Instances {
				if gceInstance.Status != "RUNNING" { // TODO: TBD, include STOPPING and SUSPENDING
					continue
				}

				instanceList = append(instanceList, instanceInfo{region: region, name: gceInstance.Name})
			}
		}

		if c.done() || ial.NextPageToken == "" {
			break
		}
		nextPageToken = ial.NextPageToken
	}
	client.CloseIdleConnections()

	if c.done() {
		return emptyList, nil
	}

	return instanceList, nil
}
