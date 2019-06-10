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

// handle AWS/EBS specific tasks
// https://docs.aws.amazon.com/elasticloadbalancing/latest/classic/EBS-cloudwatch-metrics.html

// EBS defines the collector instance
type EBS struct {
	common
}

func newEBS(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/EBS"
	c := &EBS{
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
func (c *EBS) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric{
				Name:  "VolumeReadBytes",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
				Units: "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "VolumeWriteBytes",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
				Units: "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "VolumeReadOps",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
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
				Name:  "VolumeWriteOps",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
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
				Name:  "VolumeTotalReadTime",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
				Units: "Seconds",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically

			},
		},
		{
			AWSMetric{
				Name:  "VolumeTotalWriteTime",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
				Units: "Seconds",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "VolumeIdleTime",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
				Units: "Seconds",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "VolumeQueueLength",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
				Units: "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		//
		// Use with provisioned iops ssd volumns only
		//
		// {
		// 	AWSMetric{
		// 		Name:  "VolumeThroughputPercentage",
		// 		Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
		// 		Units: "Percent",
		// 	},
		// 	CirconusMetric{
		// 		Name: "",              // NOTE: AWSMetric.Name will be used if blank
		// 		Type: "gauge",         // (gauge|counter|histogram|text)
		// 		Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
		// 	},
		// },
		// {
		// 	AWSMetric{
		// 		Name:  "VolumeConsumedReadWriteOps",
		// 		Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
		// 		Units: metricStatSum,
		// 	},
		// 	CirconusMetric{
		// 		Name: "",              // NOTE: AWSMetric.Name will be used if blank
		// 		Type: "gauge",         // (gauge|counter|histogram|text)
		// 		Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
		// 	},
		// },
		//
		// Use with general purpose ssd (gp2), throughput optimized hdd (st1), and cold hdd (sc1) volumes only
		//
		// {
		// 	AWSMetric{
		// 		Name:  "BurstBalance",
		// 		Stats: []string{metricStatAverage},
		// 		Units: "Percent",
		// 	},
		// 	CirconusMetric{
		// 		Name: "",              // NOTE: AWSMetric.Name will be used if blank
		// 		Type: "gauge",         // (gauge|counter|histogram|text)
		// 		Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
		// 	},
		// },
	}
}
