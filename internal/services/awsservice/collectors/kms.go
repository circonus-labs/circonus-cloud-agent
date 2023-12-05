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

// handle AWS/KMS specific tasks
// https://docs.aws.amazon.com/kms/latest/developerguide/monitoring-cloudwatch.html

// KMS defines the collector instance.
type KMS struct {
	common
}

func newKMS(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/KMS"
	c := &KMS{
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
// Metric is only valid for EXTERNAL CMKs
// Dimension: KeyId.
func (c *KMS) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric: AWSMetric{
				Disabled: false,
				Name:     "SecondsUntilKeyMaterialExpiration",
				Stats:    []string{metricStatMinimum},
				Units:    "Seconds",
			},
			CirconusMetric: CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
	}
}
