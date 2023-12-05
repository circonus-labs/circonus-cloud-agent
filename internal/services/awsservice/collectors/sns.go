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

// handle AWS/SNS specific tasks
// https://docs.aws.amazon.com/sns/latest/dg/sns-monitoring-using-cloudwatch.html

// SNS defines the collector instance.
type SNS struct {
	common
}

func newSNS(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/SNS"
	c := &SNS{
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
func (c *SNS) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric: AWSMetric{
				Name:  "NumberOfMessagesPublished",
				Stats: []string{metricStatSum},
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
				Name:  "NumberOfNotificationsDelivered",
				Stats: []string{metricStatSum},
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
				Name:  "NumberOfNotificationsFailed",
				Stats: []string{metricStatSum, metricStatAverage},
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
				Name:  "NumberOfNotificationsFilteredOut",
				Stats: []string{metricStatSum, metricStatAverage},
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
				Name:  "NumberOfNotificationsFilteredOutNoMessageAttributes",
				Stats: []string{metricStatSum, metricStatAverage},
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
				Name:  "NumberOfNotificationsFilteredOutInvalidAttributes",
				Stats: []string{metricStatSum, metricStatAverage},
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
				Name:  "PublishSize",
				Stats: []string{"Minimum", metricStatMaximum, metricStatAverage, metricStatSampleCount},
				Units: "Bytes",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "SMSSuccessRate",
				Stats: []string{metricStatSum, metricStatAverage, metricStatSampleCount},
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
