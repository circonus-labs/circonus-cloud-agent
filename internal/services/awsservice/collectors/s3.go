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

// handle AWS/S3 specific tasks
// https://docs.aws.amazon.com/AmazonS3/latest/dev/cloudwatch-monitoring.html

// S3 defines the collector instance.
type S3 struct {
	common
}

func newS3(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/S3"
	c := &S3{
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
func (c *S3) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "BucketSizeBytes",
				Stats:    []string{metricStatAverage},
				Units:    "Bytes",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "NumberOfObjects",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "AllRequests",
				Stats:    []string{metricStatMinimum},
				Units:    "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "GetRequests",
				Stats:    []string{metricStatMinimum},
				Units:    "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "PutRequests",
				Stats:    []string{metricStatMinimum},
				Units:    "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically

			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "DeleteRequests",
				Stats:    []string{metricStatMinimum},
				Units:    "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "HeadRequests",
				Stats:    []string{metricStatMinimum},
				Units:    "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "PostRequests",
				Stats:    []string{metricStatMinimum},
				Units:    "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "SelectRequests",
				Stats:    []string{metricStatMinimum},
				Units:    "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "SelectScannedBytes",
				Stats:    []string{metricStatAverage, metricStatMinimum, metricStatSampleCount, metricStatMinimum, metricStatMaximum},
				Units:    "Bytes",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "SelectReturnedBytes",
				Stats:    []string{metricStatAverage, metricStatMinimum, metricStatSampleCount, metricStatMinimum, metricStatMaximum},
				Units:    "Bytes",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "ListRequests",
				Stats:    []string{metricStatMinimum},
				Units:    "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "BytesDownloaded",
				Stats:    []string{metricStatAverage, metricStatAverage, metricStatSampleCount, metricStatMinimum, metricStatMaximum},
				Units:    "Bytes",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "BytesUploaded",
				Stats:    []string{metricStatAverage, metricStatMinimum, metricStatSampleCount, metricStatMinimum, metricStatMaximum},
				Units:    "Bytes",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "4xxErrors",
				Stats:    []string{metricStatAverage, metricStatMinimum, metricStatSampleCount, metricStatMinimum, metricStatMaximum},
				Units:    "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "5xxErrors",
				Stats:    []string{metricStatAverage, metricStatMinimum, metricStatSampleCount, metricStatMinimum, metricStatMaximum},
				Units:    "Count",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "FirstByteLatency",
				Stats:    []string{metricStatAverage, metricStatMinimum, metricStatSampleCount, metricStatMinimum, metricStatMaximum},
				Units:    "Milliseconds",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "TotalRequestLatency",
				Stats:    []string{metricStatAverage, metricStatMinimum, metricStatSampleCount, metricStatMinimum, metricStatMaximum},
				Units:    "Milliseconds",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
	}
}
