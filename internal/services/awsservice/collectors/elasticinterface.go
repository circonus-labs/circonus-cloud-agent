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

// handle AWS/ElasticInterface specific tasks
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/cloudwatch-metrics-ei.html#ei-metrics-dimensions

// ElasticInterface defines the collector instance.
type ElasticInterface struct {
	common
}

func newElasticInterface(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/ElasticInterface"
	c := &ElasticInterface{
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
func (c *ElasticInterface) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric{
				Name:  "AcceleratorHealthCheckFailed",
				Stats: []string{metricStatAverage},
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
				Name:  "ConnectivityCheckFailed",
				Stats: []string{metricStatAverage},
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
				Name:  "AcceleratorMemoryUsage",
				Stats: []string{metricStatSum, metricStatMinimum, metricStatMaximum, metricStatAverage},
				Units: "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
	}
}
