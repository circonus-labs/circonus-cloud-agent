// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"context"

	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/rs/zerolog"
)

// handle AWS/Route53Resolver specific tasks
// https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/monitoring-resolver-with-cloudwatch.html

// Route53Resolver defines the collector instance.
type Route53Resolver struct {
	common
}

func newRoute53Resolver(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/Route53Resolver"
	c := &Route53Resolver{
		common: newCommon(ctx, ns, check, cfg, logger),
	}
	if len(c.metrics) == 0 {
		c.metrics = c.DefaultMetrics()
	}
	c.tags = append(c.tags, circonus.Tag{Category: "service", Value: ns})
	c.logger.Debug().Msg("initialized")
	return c, nil
}

// DefaultMetrics returns a default metric configuration.
func (c *Route53Resolver) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric{
				Name:  "InboundQueryVolume",
				Stats: []string{metricStatSum},
				Units: "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "OutboundQueryVolume",
				Stats: []string{metricStatSum},
				Units: "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
	}
}
