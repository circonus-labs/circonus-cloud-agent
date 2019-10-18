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

// handle AWS/DX specific tasks
// https://docs.aws.amazon.com/directconnect/latest/UserGuide/monitoring-cloudwatch.html

// DX defines the collector instance
type DX struct {
	common
}

func newDX(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/DX"
	c := &DX{
		common: newCommon(ctx, ns, check, cfg, logger),
	}
	if len(c.metrics) == 0 {
		c.metrics = c.DefaultMetrics()
	}
	c.tags = append(c.tags, circonus.Tag{Category: "service", Value: ns})
	c.logger.Debug().Msg("initialized")
	return c, nil
}

// DefaultMetrics returns a default metric configuration
func (c *DX) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric{
				Name:  "ConnectionState",
				Stats: []string{metricStatMinimum, metricStatMaximum},
				Units: "Status",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "text",          // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "ConnectionBpsEgress",
				Stats: []string{metricStatMinimum, metricStatAverage, metricStatMaximum},
				Units: "Bits",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "ConnectionBpsIngress",
				Stats: []string{metricStatMinimum, metricStatAverage, metricStatMaximum},
				Units: "Bits",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "ConnectionPpsEgress",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage},
				Units: "Packets",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "ConnectionPpsIngress",
				Stats: []string{metricStatMaximum, metricStatAverage, metricStatMinimum},
				Units: "Packets",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically

			},
		},
		{
			AWSMetric{
				Name:  "ConnectionCRCErrorCount",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum},
				Units: "Errors",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "ConnectionLightLevelTx",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage},
				Units: "dBm",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "ConnectionLightLevelRx",
				Stats: []string{metricStatMaximum, metricStatAverage, metricStatMinimum},
				Units: "dBm",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
	}
}
