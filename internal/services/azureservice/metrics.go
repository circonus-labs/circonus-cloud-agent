// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package azureservice

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights"
	"github.com/Azure/go-autorest/autorest"
	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/pkg/errors"
)

// getResourceMetrics collects metrics using azure api for a given resource id
// and writes them to the metric destination
func (inst *Instance) getResourceMetrics(
	metricDest io.Writer,
	auth autorest.Authorizer,
	resourceID string,
	endTime time.Time,
	resourceTags circonus.Tags) error {

	// list keyed on granularity (PT1M, PT5M, ...)-> aggregation (Average, Count, Maximum, Minimum, Total) -> list of metric names
	// for clarity and reference: map[granularity]map[aggregation][]metric_name
	metricList, err := inst.getMetricList(auth, resourceID)
	if err != nil {
		return err
	}

	for granularity, aggregations := range metricList {
		for aggregation, metrics := range aggregations {

			if inst.done() {
				return nil
			}

			// limit on number of metrics that can be requested is 20
			// split larger metric list into groups of 20 metrics each
			//
			// limit is NOT documented here: https://godoc.org/github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights#MetricsClient.List
			// but List request returns: "insights.MetricsClient#List: Failure responding to request: StatusCode=400 -- Original Error: autorest/azure: Service returned an error. Status=400 Code=\"BadRequest\" Message=\"Requested metrics count: 22 bigger than allowed max: 20\""
			var metricGroups [][]string
			maxGroupSize := 20
			for i := 0; i < len(metrics); i += maxGroupSize {
				end := i + maxGroupSize
				if end > len(metrics) {
					end = len(metrics)
				}
				metricGroups = append(metricGroups, metrics[i:end])
			}

			for _, metricGroup := range metricGroups {
				err := inst.handleMetricGroup(metricDest, auth, resourceID, resourceTags, endTime, granularity, aggregation, metricGroup)
				if err != nil {
					inst.logger.Error().
						Err(err).
						Str("resource_id", resourceID).
						Str("granularity", granularity).
						Str("aggregation", aggregation).
						Strs("metrics", metricGroup).
						Msg("handling metric samples")
				}
			}
		}
	}

	return nil
}

// handleMetricGroup retrieves metric samples for a group of metrics (based on
// azure api limit) and handles processing and writing each metric sample to the
// metric destination
func (inst *Instance) handleMetricGroup(
	metricDest io.Writer,
	auth autorest.Authorizer,
	resourceID string,
	resourceTags circonus.Tags,
	endTime time.Time,
	granularity string,
	aggregation string,
	metricGroup []string) error {

	metricData, err := inst.getMetricData(auth, resourceID, endTime, granularity, aggregation, metricGroup)
	if err != nil {
		return errors.Wrap(err, "fetching metric samples")
	}

	for metricName, metricSamples := range metricData {
		if inst.done() {
			return nil
		}

		var tags circonus.Tags

		tags = append(tags, inst.baseTags...)
		tags = append(tags, resourceTags...)
		tags = append(tags, circonus.Tags{
			circonus.Tag{Category: "units", Value: metricSamples.Units},
			circonus.Tag{Category: "aggregation", Value: aggregation}}...)

		encodedMetricName := inst.check.MetricNameWithStreamTags(metricName, tags)

		for _, sample := range metricSamples.Samples {
			if inst.done() {
				return nil
			}

			err := inst.check.WriteMetricSample(metricDest, encodedMetricName, sample.Type, sample.Value, &sample.Timestamp)
			if err != nil {
				inst.logger.Warn().
					Err(err).
					Str("metric_name", metricName).
					Interface("value", sample.Value).
					Str("type", sample.Type).
					Time("ts", sample.Timestamp).
					Str("encoded_metric_name", encodedMetricName).
					Interface("tags", tags).
					Msg("adding metric sample")
			}
		}
	}

	return nil
}

type metricSample struct {
	Timestamp time.Time
	Value     interface{}
	Type      string
}
type metricSamples struct {
	Units   string
	Samples []metricSample
}
type metricData map[string]metricSamples

// getMetricData uses Azure API to retrieve metric samples for a resource
func (inst *Instance) getMetricData(
	auth autorest.Authorizer,
	resourceID string,
	endTime time.Time,
	granularity string,
	aggregation string,
	metricList []string) (metricData, error) {

	if len(metricList) == 0 {
		return metricData{}, nil
	}

	metricsClient := insights.NewMetricsClient(inst.cfg.Azure.SubscriptionID)
	metricsClient.Authorizer = auth
	if err := metricsClient.AddToUserAgent(inst.cfg.Azure.UserAgent); err != nil {
		inst.logger.Warn().Err(err).Msg("adding user agent to client")
	}

	timeUnit := time.Minute
	switch granularity {
	case "PT5M":
		timeUnit = 5 * time.Minute
	case "PT15M":
		timeUnit = 15 * time.Minute
	case "PT30M":
		timeUnit = 30 * time.Minute
	case "PT1H":
		timeUnit = time.Hour
	}

	timeDelta := timeUnit * 10 // get last 10 samples

	startTime := endTime.Add(-timeDelta)
	timespan := fmt.Sprintf("%s/%s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	resp, err := metricsClient.List(inst.ctx, resourceID, timespan, &granularity, strings.Join(metricList, ","), aggregation, nil, "", "", insights.Data, "")
	if err != nil {
		return nil, err
	}

	data := metricData{}
	for _, v := range *resp.Value {
		metricName := *v.Name.Value
		metricUnits := string(v.Unit)
		samples := inst.extractSamples(v.Timeseries, aggregation)
		if len(samples) > 0 {
			data[metricName] = metricSamples{Units: metricUnits, Samples: samples}
		}
	}

	return data, nil
}

// extractSamples takes the azure metric timeseries data and returns a list where
// each item in the list is a metric sample comprised of a timestamp and the sample's value.
func (inst *Instance) extractSamples(timeseries *[]insights.TimeSeriesElement, aggregation string) []metricSample {
	var samples []metricSample
	for _, t := range *timeseries {
		// ref: https://godoc.org/github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights#MetricValue
		for _, mv := range *t.Data {
			var sampleValue interface{}
			sampleType := circonus.MetricTypeFloat64
			sampleTimestamp := mv.TimeStamp.ToTime()

			switch aggregation {
			case "Average":
				if mv.Average != nil {
					sampleValue = *mv.Average
				}
			case "Count":
				if mv.Count != nil {
					sampleValue = *mv.Count
					sampleType = circonus.MetricTypeUint64
				}
			case "Maximum":
				if mv.Maximum != nil {
					sampleValue = *mv.Maximum
				}
			case "Minimum":
				if mv.Minimum != nil {
					sampleValue = *mv.Minimum
				}
			case "Total":
				if mv.Total != nil {
					sampleValue = *mv.Total
				}
			}

			if sampleValue != nil {
				samples = append(samples, metricSample{Timestamp: sampleTimestamp, Type: sampleType, Value: sampleValue})
			}
		}
	}
	return samples
}

// getMetricList retrieves a list of viable metrics for a given resource
func (inst *Instance) getMetricList(auth autorest.Authorizer, resourceID string) (map[string]map[string][]string, error) {
	// ref: https://godoc.org/github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights

	// get metric definitions for resource
	metricsDefClient := insights.NewMetricDefinitionsClient(inst.cfg.Azure.SubscriptionID)
	metricsDefClient.Authorizer = auth
	if err := metricsDefClient.AddToUserAgent(inst.cfg.Azure.UserAgent); err != nil {
		inst.logger.Warn().Err(err).Msg("adding user agent to client")
	}
	// ref: https://godoc.org/github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights#MetricDefinitionsClient.List
	result, err := metricsDefClient.List(inst.ctx, resourceID, "")
	if err != nil {
		return nil, err
	}

	// granularities by preference
	granularities := []string{"PT1M", "PT5M", "PT15M", "PT30M", "PT1H"}

	// build lists based on PrimaryAggregationType for each metric - the primary aggregation is azure's suggested use/exposure
	// ref: https://godoc.org/github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights#AggregationType
	ml := map[string]map[string][]string{}
	for i := range *result.Value {
		// ref: https://godoc.org/github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights#MetricDefinition
		if strings.Contains(*(*result.Value)[i].Name.LocalizedValue, "(Deprecated)") {
			continue
		}
		if strings.Contains(*(*result.Value)[i].Name.LocalizedValue, "(Preview)") {
			continue
		}

		if (*result.Value)[i].PrimaryAggregationType == insights.None {
			// skip 'None' aggregation type; there's no 'None' value
			// in MetricValue struct, nor any documentation on what
			// would be returned
			continue
		}

		// find smallest granularity the metric provides
		granularity := ""
		for _, g := range granularities {
			for _, ma := range *(*result.Value)[i].MetricAvailabilities {
				if *ma.TimeGrain == g {
					granularity = g
					break
				}
			}
			if granularity != "" {
				break
			}
		}

		metricName := *(*result.Value)[i].Name.Value
		metricAggregation := string((*result.Value)[i].PrimaryAggregationType)
		if _, found := ml[granularity]; !found {
			ml[granularity] = map[string][]string{}
		}
		if _, found := ml[granularity][metricAggregation]; !found {
			ml[granularity][metricAggregation] = []string{}
		}

		ml[granularity][metricAggregation] = append(ml[granularity][metricAggregation], metricName)
	}

	return ml, nil
}
