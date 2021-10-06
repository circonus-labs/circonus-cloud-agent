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

// handle AWS/RDS specific tasks
// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/MonitoringOverview.html#monitoring-cloudwatch

// RDS defines the collector instance.
type RDS struct {
	common
}

func newRDS(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/RDS"
	c := &RDS{
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
func (c *RDS) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric{
				Disabled: false,
				Name:     "BinLogDiskUsage",
				Stats:    []string{metricStatAverage},
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
				Name:     "BurstBalance",
				Stats:    []string{metricStatAverage},
				Units:    "Percent",
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
				Name:     "CPUUtilization",
				Stats:    []string{metricStatAverage},
				Units:    "Percent",
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
				Name:     "CPUCreditUsage",
				Stats:    []string{metricStatAverage},
				Units:    "Credits",
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
				Name:     "CPUCreditBalance",
				Stats:    []string{metricStatAverage},
				Units:    "Credits",
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
				Name:     "DatabaseConnections",
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
				Name:     "DiskQueueDepth",
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
				Name:     "FailedSQLServerAgentJobsCount",
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
				Name:     "FreeableMemory",
				Stats:    []string{metricStatAverage},
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
				Name:     "FreeStorageSpace",
				Stats:    []string{metricStatAverage},
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
				Name:     "MaximumUsedTransactionIDs",
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
				Name:     "NetworkReceiveThroughput",
				Stats:    []string{metricStatAverage},
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
				Name:     "NetworkTransmitThroughput",
				Stats:    []string{metricStatAverage},
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
				Name:     "OldestReplicationSlotLag",
				Stats:    []string{metricStatAverage},
				Units:    "Megabytes",
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
				Name:     "ReadIOPS",
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
				Name:     "ReadLatency",
				Stats:    []string{metricStatAverage},
				Units:    "Seconds",
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
				Name:     "ReadThroughput",
				Stats:    []string{metricStatAverage},
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
				Name:     "ReplicaLag",
				Stats:    []string{metricStatAverage},
				Units:    "Seconds",
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
				Name:     "ReplicationSlotDiskUsage",
				Stats:    []string{metricStatAverage},
				Units:    "Megabytes",
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
				Name:     "SwapUsage",
				Stats:    []string{metricStatAverage},
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
				Name:     "TransactionLogDiskUsage",
				Stats:    []string{metricStatAverage},
				Units:    "Megabytes",
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
				Name:     "TransactionLogsGeneration",
				Stats:    []string{metricStatAverage},
				Units:    "Megabytes",
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
				Name:     "WriteIOPS",
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
				Name:     "WriteLatency",
				Stats:    []string{metricStatAverage},
				Units:    "Seconds",
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
				Name:     "WriteThroughput",
				Stats:    []string{metricStatAverage},
				Units:    "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
	}
}
