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

// handle AWS/EC2Spot specific tasks
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/spot-fleet-cloudwatch-metrics.html

// EC2Spot defines the collector instance.
type EC2Spot struct {
	common
}

func newEC2Spot(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/EC2Spot"
	c := &EC2Spot{
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
func (c *EC2Spot) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric: AWSMetric{
				Name:  "AvailableInstancePoolsCount",
				Stats: []string{metricStatAverage},
				Units: "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "BidsSubmittedForCapacity",
				Stats: []string{metricStatAverage},
				Units: "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "EligibleInstancePoolCount",
				Stats: []string{metricStatAverage},
				Units: "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "FulfilledCapacity",
				Stats: []string{metricStatAverage},
				Units: "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "MaxPercentCapacityAllocation",
				Stats: []string{metricStatAverage},
				Units: "Percent",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "PendingCapacity",
				Stats: []string{metricStatAverage},
				Units: "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically

			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "PercentCapacityAllocation",
				Stats: []string{metricStatAverage},
				Units: "Percent",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "TargetCapacity",
				Stats: []string{metricStatAverage},
				Units: "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "TerminatingCapacity",
				Stats: []string{metricStatAverage},
				Units: "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
	}
}
