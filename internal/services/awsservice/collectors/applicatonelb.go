// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"context"
	"strings"

	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// handle AWS/ApplicationELB specific tasks
// https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-cloudwatch-metrics.html

// ApplicationELB defines the collector instance.
type ApplicationELB struct {
	common
}

func newApplicationELB(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	if len(cfg.Dimensions) == 0 {
		return nil, errors.New("metrics *require* dimension(s)")
	}
	ns := "AWS/ApplicationELB"
	c := &ApplicationELB{
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
func (c *ApplicationELB) DefaultMetrics() []Metric {
	haveLoadBalancer := false
	haveAvailabilityZone := false
	haveTargetGroup := false
	for _, dim := range c.dimensions {
		switch strings.ToLower(*dim.Name) {
		case "loadbalancer":
			haveLoadBalancer = true
		case "availabilityzone":
			haveAvailabilityZone = true
		case "targetgroup":
			haveTargetGroup = true
		}
	}

	if haveTargetGroup && haveAvailabilityZone && haveLoadBalancer {
		return []Metric{
			{
				AWSMetric: AWSMetric{
					Name:  "IPv6RequestCount",
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
			{
				AWSMetric: AWSMetric{
					Name:  "HealthyHostCount",
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
					Name:  "HTTPCode_Target_2XX_Count",
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
					Name:  "HTTPCode_Target_3XX_Count",
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
					Name:  "HTTPCode_Target_4XX_Count",
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
					Name:  "HTTPCode_Target_5XX_Count",
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
					Name:  "NonStickyRequestCount",
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
					Name:  "TargetConnectionErrorCount",
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
					Name:  "TargetResponseTime",
					Stats: []string{metricStatSum},
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
					Name:  "TargetTLSNegotiationErrorCount",
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
					Name:  "UnHealthyHostCount",
					Stats: []string{metricStatAverage, metricStatMinimum, metricStatMaximum},
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

	if haveTargetGroup && haveLoadBalancer {
		return []Metric{
			{
				AWSMetric: AWSMetric{
					Name:  "IPv6RequestCount",
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
			{
				AWSMetric: AWSMetric{
					Name:  "HealthyHostCount",
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
					Name:  "HTTPCode_Target_2XX_Count",
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
					Name:  "HTTPCode_Target_3XX_Count",
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
					Name:  "HTTPCode_Target_4XX_Count",
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
					Name:  "HTTPCode_Target_5XX_Count",
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
					Name:  "NonStickyRequestCount",
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
					Name:  "RequestCountPerTarget",
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
					Name:  "TargetConnectionErrorCount",
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
					Name:  "TargetResponseTime",
					Stats: []string{metricStatSum},
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
					Name:  "TargetTLSNegotiationErrorCount",
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
					Name:  "UnHealthyHostCount",
					Stats: []string{metricStatAverage, metricStatMinimum, metricStatMaximum},
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

	if haveAvailabilityZone && haveLoadBalancer {
		return []Metric{
			{
				AWSMetric: AWSMetric{
					Name:  "ClientTLSNegotiationErrorCount",
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
					Name:  "HTTPCode_ELB_4XX_Count",
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
					Name:  "HTTPCode_ELB_5XX_Count",
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
					Name:  "IPv6RequestCount",
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
					Name:  "RejectedConnectionCount",
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
			{
				AWSMetric: AWSMetric{
					Name:  "HTTPCode_Target_2XX_Count",
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
					Name:  "HTTPCode_Target_3XX_Count",
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
					Name:  "HTTPCode_Target_4XX_Count",
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
					Name:  "HTTPCode_Target_5XX_Count",
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
					Name:  "TargetConnectionErrorCount",
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
					Name:  "TargetResponseTime",
					Stats: []string{metricStatSum},
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
					Name:  "TargetTLSNegotiationErrorCount",
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

	if haveTargetGroup {
		return []Metric{
			{
				AWSMetric: AWSMetric{
					Name:  "NonStickyRequestCount",
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
					Name:  "RequestCountPerTarget",
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
					Name:  "LambdaInternalError",
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
					Name:  "LambdaUserError",
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

	if haveLoadBalancer {
		return []Metric{
			{
				AWSMetric: AWSMetric{
					Name:  "ActiveConnectionCount",
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
					Name:  "ClientTLSNegotiationErrorCount",
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
					Name:  "ConsumedLCUs",
					Stats: []string{"Minimum", "Maximum", "Average", metricStatSum},
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
					Name:  "HTTP_Fixed_Response_Count",
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
					Name:  "HTTP_Redirect_Count",
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
					Name:  "HTTP_Redirect_Url_Limit_Exceeded_Count",
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
					Name:  "HTTPCode_ELB_3XX_Count",
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
					Name:  "HTTPCode_ELB_4XX_Count",
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
					Name:  "HTTPCode_ELB_5XX_Count",
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
					Name:  "IPv6ProcessedBytes",
					Stats: []string{metricStatSum},
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
					Name:  "IPv6RequestCount",
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
					Name:  "NewConnectionCount",
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
					Name:  "ProcessedBytes",
					Stats: []string{metricStatSum},
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
					Name:  "RejectedConnectionCount",
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
			{
				AWSMetric: AWSMetric{
					Name:  "RuleEvaluations",
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
					Name:  "StandardProcessedBytes",
					Stats: []string{metricStatSum},
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
					Name:  "HTTPCode_Target_2XX_Count",
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
					Name:  "HTTPCode_Target_3XX_Count",
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
					Name:  "HTTPCode_Target_4XX_Count",
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
					Name:  "HTTPCode_Target_5XX_Count",
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
					Name:  "TargetConnectionErrorCount",
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
					Name:  "TargetResponseTime",
					Stats: []string{metricStatSum},
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
					Name:  "TargetTLSNegotiationErrorCount",
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
					Name:  "LambdaInternalError",
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
					Name:  "LambdaTargetProcessedBytes",
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
					Name:  "LambdaUserError",
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
					Name:  "ELBAuthError",
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
					Name:  "ELBAuthFailure",
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
					Name:  "ELBAuthLatency",
					Stats: []string{metricStatAverage, metricStatMinimum, metricStatMaximum, metricStatSum, metricStatSampleCount},
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
					Name:  "ELBAuthRefreshTokenSuccess",
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
					Name:  "ELBAuthSuccess",
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
					Name:  "ELBAuthUserClaimsSizeExceeded",
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

	return []Metric{}
}
