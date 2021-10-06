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

// handle AWS/Lambda specific tasks
// https://docs.aws.amazon.com/vpc/latest/tgw/transit-gateway-cloudwatch-metrics.html

// Lambda defines the collector instance.
type Lambda struct {
	common
}

func newLambda(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/Lambda"
	c := &Lambda{
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
func (c *Lambda) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric{
				Disabled: false,
				Name:     "Invocations",
				Stats:    []string{metricStatAverage, metricStatSum},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: false,
				Name:     "Errors",
				Stats:    []string{metricStatAverage, metricStatSum},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: false,
				Name:     "DeadLetterErrors",
				Stats:    []string{metricStatAverage, metricStatSum},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: false,
				Name:     "Duration",
				Stats:    []string{metricStatAverage, metricStatSum},
				Units:    "Milliseconds",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: false,
				Name:     "Throttles",
				Stats:    []string{metricStatAverage, metricStatSum},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically

			},
		},
		{
			AWSMetric{
				Disabled: false,
				Name:     "IteratorAge",
				Stats:    []string{metricStatAverage},
				Units:    "Milliseconds",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: false,
				Name:     "ConcurrentExecutions",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: false,
				Name:     "UnreservedConcurrentExecutions",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
	}
}
