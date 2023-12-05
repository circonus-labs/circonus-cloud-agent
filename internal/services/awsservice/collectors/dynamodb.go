// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"context"

	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// handle AWS/DynamoDB specific tasks
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/metrics-dimensions.html
// NOTE: the dynamo db metrics pretty much require user defines dimensions - see link for details
//       regarding which dimensions apply to which metrics

// DynamoDB defines the collector instance.
type DynamoDB struct {
	common
}

func newDynamoDB(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	if len(cfg.Dimensions) == 0 {
		return nil, errors.New("metrics *require* dimensions")
	}
	ns := "AWS/DynamoDB"
	c := &DynamoDB{
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
func (c *DynamoDB) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric: AWSMetric{
				Name:  "ConditionalCheckFailedRequests",
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
				Name:  "ConsumedReadCapacityUnits",
				Stats: []string{metricStatAverage, metricStatSum},
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
				Name:  "ConsumedReadCapacityUnits",
				Stats: []string{metricStatAverage, metricStatSum},
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
				Name:  "OnlineIndexConsumedWriteCapacity",
				Stats: []string{metricStatAverage, metricStatSum},
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
				Name:  "PendingReplicationCount",
				Stats: []string{metricStatAverage, metricStatSum},
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
				Name:  "ProvisionedReadCapacityUnits",
				Stats: []string{metricStatAverage, metricStatMinimum, metricStatMaximum},
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
				Name:  "ProvisionedWriteCapacityUnits",
				Stats: []string{metricStatAverage, metricStatMinimum, metricStatMaximum},
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
				Name:  "ReadThrottleEvents",
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
				Name:  "ReplicationLatency",
				Stats: []string{metricStatAverage, metricStatMinimum, metricStatMaximum},
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
				Name:  "ReturnedBytes",
				Stats: []string{metricStatAverage, metricStatMinimum, metricStatMaximum, metricStatSum},
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
				Name:  "ReturnedItemCount",
				Stats: []string{metricStatAverage, metricStatMinimum, metricStatMaximum, metricStatSum},
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
				Name:  "ReturnedRecordsCount",
				Stats: []string{metricStatAverage, metricStatMinimum, metricStatMaximum, metricStatSum},
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
				Name:  "SuccessfulRequestLatency",
				Stats: []string{metricStatAverage, metricStatMinimum, metricStatMaximum, metricStatSum},
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
				Name:  "SystemErrors",
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
				Name:  "TimeToLiveDeletedItemCount",
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
				Name:  "ThrottledRequests",
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
				Name:  "UserErrors",
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
				Name:  "WriteThrottleEvents",
				Stats: []string{metricStatSum},
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
