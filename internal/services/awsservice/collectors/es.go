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

// handle AWS/ES specific tasks
// https://docs.aws.amazon.com/elasticsearch-service/latest/developerguide/es-managedomains.html#es-managedomains-cloudwatchmetrics

// ES defines the collector instance.
type ES struct {
	common
}

func newES(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/ES"
	c := &ES{
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
func (c *ES) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric: AWSMetric{
				Name:  "ClusterStatus.green",
				Stats: []string{metricStatMinimum, metricStatMaximum},
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
				Name:  "ClusterStatus.yellow",
				Stats: []string{metricStatMinimum, metricStatMaximum},
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
				Name:  "ClusterStatus.red",
				Stats: []string{metricStatMinimum, metricStatMaximum},
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
				Name:  "Nodes",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage},
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
				Name:  "CPUUtilization",
				Stats: []string{metricStatMaximum, metricStatAverage},
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
				Name:  "FreeStorageSpace",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum},
				Units: "GiB",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "ClusterUsedSpace",
				Stats: []string{metricStatMinimum, metricStatMaximum},
				Units: "GiB",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "ClusterIndexWritesBlocked",
				Stats: []string{metricStatMaximum},
				Units: "Status",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "JVMMemoryPressure",
				Stats: []string{metricStatMaximum},
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
				Name:  "AutomatedSnapshotFailure",
				Stats: []string{metricStatMinimum, metricStatMaximum},
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
				Name:  "CPUCreditBalance",
				Stats: []string{metricStatMinimum},
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
				Name:  "KibanaHealthyNodes",
				Stats: []string{metricStatMinimum},
				Units: "Status",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "KMSKeyError",
				Stats: []string{metricStatMinimum, metricStatMaximum},
				Units: "Status",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "KMSKeyInaccessible",
				Stats: []string{metricStatMinimum, metricStatMaximum},
				Units: "Status",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "InvalidHostHeaderRequests",
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
				Name:  "ElasticsearchRequests",
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
				Name:  "RequestCount",
				Stats: []string{metricStatSum},
				Units: "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		// dedicated master node metrics (but the doc doesn't indicate that aws/es takes a "node" dimension???)
		{
			AWSMetric: AWSMetric{
				Name:  "MasterCPUUtilization",
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
				Name:  "MasterJVMMemoryPressure",
				Stats: []string{metricStatMaximum},
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
				Name:  "MasterCPUCreditBalance",
				Stats: []string{metricStatMinimum},
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
				Name:  "MasterReachableFromNode",
				Stats: []string{metricStatMinimum},
				Units: "Status",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		// EBS volume metrics (but the doc doesn't indicate that aws/es takes an "ebs volume" dimension???)
		{
			AWSMetric: AWSMetric{
				Name:  "ReadLatency",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage},
				Units: "Seconds",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "WriteLatency",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage},
				Units: "Seconds",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "ReadThroughput",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage},
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
				Name:  "WriteThroughput",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage},
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
				Name:  "DiskQueueDepth",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage},
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
				Name:  "ReadIOPS",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage},
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
				Name:  "WriteIOPS",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage},
				Units: "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		// instance metrics (but the doc doesn't indicate that aws/es takes an "instance" dimension???)
		{
			AWSMetric: AWSMetric{
				Name:  "IndexingLatency",
				Stats: []string{metricStatAverage},
				Units: "Milliseconds",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "IndexingRate",
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
				Name:  "SearchLatency",
				Stats: []string{metricStatAverage},
				Units: "Milliseconds",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "SearchRate",
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
				Name:  "SysMemoryUtilization",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage},
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
				Name:  "JVMGCYoungCollectionCount",
				Stats: []string{metricStatMaximum},
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
				Name:  "JVMGCYoungCollectionTime",
				Stats: []string{metricStatMaximum},
				Units: "Milliseconds",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "JVMGCOldCollectionCount",
				Stats: []string{metricStatMaximum},
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
				Name:  "JVMGCOldCollectionTime",
				Stats: []string{metricStatMaximum},
				Units: "Milliseconds",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Name:  "ThreadpoolForce_mergeQueue",
				Stats: []string{metricStatMaximum},
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
				Name:  "ThreadpoolForce_mergeRejected",
				Stats: []string{metricStatMaximum},
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
				Name:  "ThreadpoolForce_mergeThreads",
				Stats: []string{metricStatMaximum},
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
				Name:  "ThreadpoolIndexQueue",
				Stats: []string{metricStatMaximum},
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
				Name:  "ThreadpoolIndexRejected",
				Stats: []string{metricStatMaximum},
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
				Name:  "ThreadpoolIndexThreads",
				Stats: []string{metricStatMaximum},
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
				Name:  "ThreadpoolSearchQueue",
				Stats: []string{metricStatMaximum},
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
				Name:  "ThreadpoolSearchRejected",
				Stats: []string{metricStatMaximum},
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
				Name:  "ThreadpoolSearchThreads",
				Stats: []string{metricStatMaximum},
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
				Name:  "ThreadpoolBulkQueue",
				Stats: []string{metricStatMaximum},
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
				Name:  "ThreadpoolBulkRejected",
				Stats: []string{metricStatMaximum},
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
				Name:  "ThreadpoolBulkThreads",
				Stats: []string{metricStatMaximum},
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
