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

// handle AWS/TransitGateway specific tasks
// https://docs.aws.amazon.com/vpc/latest/tgw/transit-gateway-cloudwatch-metrics.html

// TransitGateway defines the collector instance.
type TransitGateway struct {
	common
}

func newTransitGateway(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/TransitGateway"
	c := &TransitGateway{
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
func (c *TransitGateway) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric{
				Disabled: false,
				Name:     "BytesIn",
				Stats:    []string{metricStatAverage, metricStatSum, metricStatSampleCount},
				Units:    "Bytes",
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
				Name:     "BytesOut",
				Stats:    []string{metricStatAverage, metricStatSum, metricStatSampleCount},
				Units:    "Bytes",
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
				Name:     "PacketsIn",
				Stats:    []string{metricStatAverage, metricStatSum, metricStatSampleCount},
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
				Name:     "PacketsOut",
				Stats:    []string{metricStatAverage, metricStatSum, metricStatSampleCount},
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
				Name:     "PacketDropCountBlackhole",
				Stats:    []string{metricStatAverage, metricStatSum, metricStatSampleCount},
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
				Name:     "PacketDropCountNoRoute",
				Stats:    []string{metricStatAverage, metricStatSum, metricStatSampleCount},
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
