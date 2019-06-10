// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"context"
	"strings"

	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/rs/zerolog"
)

// handle AWS/ElasticTranscoder specific tasks
// https://docs.aws.amazon.com/elastictranscoder/latest/developerguide/metrics-dimensions.html

// ElasticTranscoder defines the collector instance
type ElasticTranscoder struct {
	common
}

func newElasticTranscoder(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/ElasticTranscoder"
	c := &ElasticTranscoder{
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
func (c *ElasticTranscoder) DefaultMetrics() []Metric {
	havePipelineID := false
	haveOperation := false
	for _, dim := range c.dimensions {
		switch strings.ToLower(*dim.Name) {
		case "pipelineid":
			havePipelineID = true
		case "operation":
			haveOperation = true
		}
	}

	if !havePipelineID && !haveOperation {
		return []Metric{
			{
				AWSMetric{
					Name:  "BilledHDOutput",
					Stats: []string{metricStatAverage},
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
					Name:  "BilledSDOutput",
					Stats: []string{metricStatAverage},
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
					Name:  "BilledAudioOutput",
					Stats: []string{metricStatAverage},
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
					Name:  "JobsErrored",
					Stats: []string{metricStatAverage, metricStatSum},
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
					Name:  "OutputsPerJob",
					Stats: []string{metricStatAverage},
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
					Name:  "StandbyTime",
					Stats: []string{metricStatAverage},
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
					Name:  "Errors",
					Stats: []string{metricStatAverage, metricStatSum},
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
					Name:  "Throttles",
					Stats: []string{metricStatAverage, metricStatSum},
					Units: "Count",
				},
				CirconusMetric{
					Name: "",              // NOTE: AWSMetric.Name will be used if blank
					Type: "gauge",         // (gauge|counter|histogram|text)
					Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
				},
			},
		}
	} else if havePipelineID && !haveOperation {
		return []Metric{
			{
				AWSMetric{
					Name:  "BilledHDOutput",
					Stats: []string{metricStatAverage},
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
					Name:  "BilledSDOutput",
					Stats: []string{metricStatAverage},
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
					Name:  "BilledAudioOutput",
					Stats: []string{metricStatAverage},
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
					Name:  "JobsCompleted",
					Stats: []string{metricStatAverage, metricStatSum},
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
					Name:  "JobsErrored",
					Stats: []string{metricStatAverage, metricStatSum},
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
					Name:  "OutputsPerJob",
					Stats: []string{metricStatAverage},
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
					Name:  "StandbyTime",
					Stats: []string{metricStatAverage},
					Units: "Count",
				},
				CirconusMetric{
					Name: "",              // NOTE: AWSMetric.Name will be used if blank
					Type: "gauge",         // (gauge|counter|histogram|text)
					Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
				},
			},
		}
	} else if !havePipelineID && haveOperation {
		return []Metric{
			{
				AWSMetric{
					Name:  "Errors",
					Stats: []string{metricStatAverage, metricStatSum},
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
					Name:  "Throttles",
					Stats: []string{metricStatAverage, metricStatSum},
					Units: "Count",
				},
				CirconusMetric{
					Name: "",              // NOTE: AWSMetric.Name will be used if blank
					Type: "gauge",         // (gauge|counter|histogram|text)
					Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
				},
			},
		}
	}

	return []Metric{}
}
