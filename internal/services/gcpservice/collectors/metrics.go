// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"fmt"
	"io"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/api/metric"
)

// processMetrics retrieves the available metrics for a resource identified by the supplied filter.
func (c *common) processMetrics(projectID, filter string, creds []byte, metricDest io.Writer, baseTags circonus.Tags) error {
	client, err := monitoring.NewMetricClient(c.ctx, option.WithCredentialsJSON(creds))
	if err != nil {
		return errors.Wrap(err, "gcp monitoring client")
	}

	req := &monitoringpb.ListMetricDescriptorsRequest{
		Name:     "projects/" + projectID,
		PageSize: 10,
		Filter:   filter,
	}

	// c.logger.Debug().Str("filter", filter).Msg("getting metric descriptors")
	metricDescriptors := make([]*metric.MetricDescriptor, 0)
	iter := client.ListMetricDescriptors(c.ctx, req)
	for {
		metricDescriptor, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			c.logger.Warn().Err(err).Interface("err_val", err).Str("filter", filter).Msg("metric descriptor, iter.next, skipping remainder")
			break
		}
		metricDescriptors = append(metricDescriptors, metricDescriptor)
	}

	c.logger.Debug().Str("filter", filter).Int("metrics", len(metricDescriptors)).Msg("processing metric descriptors")
	for _, metricDescriptor := range metricDescriptors {
		var tags circonus.Tags
		tags = append(tags, baseTags...)

		unit := metricDescriptor.GetUnit()
		if unit != "" && unit != "1" {
			tags = append(tags, circonus.Tag{Category: "units", Value: unit})
		}

		tsFilter := fmt.Sprintf(`metric.type = "%s" %s`, metricDescriptor.GetType(), filter)
		c.fetchTimeseries(client, projectID, tsFilter, creds, metricDescriptor.GetDisplayName(), metricDest, tags)
		// fetchTimeseries will log its own errors and ignore them so it can get at least some metrics
		// if err := c.fetchTimeseries(client, projectID, tsFilter, creds, metricDescriptor.GetDisplayName(), metricDest, tags); err != nil {
		// 	c.logger.Warn().Err(err).Interface("metric", metricDescriptor).Str("filter", tsFilter).Msg("fetching timeseries")
		// }
		if c.done() {
			break
		}
	}

	client.Close()
	return nil
}

// fetchTimeseries retrieves the actual samples for the metric defined by the filter.
func (c *common) fetchTimeseries(client *monitoring.MetricClient, projectID, filter string, creds []byte, metricName string, metricDest io.Writer, baseTags circonus.Tags) {
	_ = creds // ref to keep signatures same and squelch lint unused warning

	req := &monitoringpb.ListTimeSeriesRequest{
		Name: "projects/" + projectID,
		Interval: &monitoringpb.TimeInterval{
			StartTime: &timestamp.Timestamp{
				Seconds: c.tsStart.UTC().Unix(),
			},
			EndTime: &timestamp.Timestamp{
				Seconds: c.tsEnd.UTC().Unix(),
			},
		},
		Filter: filter,
		View:   monitoringpb.ListTimeSeriesRequest_FULL,
	}

	// c.logger.Debug().Str("filter", filter).Msg("getting time series list")
	timeSeriesList := make([]*monitoringpb.TimeSeries, 0)
	it := client.ListTimeSeries(c.ctx, req)
	for {
		timeSeries, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			c.logger.Warn().Err(err).Interface("err_val", err).Str("filter", filter).Msg("metric timeseries, iter.next, skipping remainder")
			break
		}
		timeSeriesList = append(timeSeriesList, timeSeries)
	}

	// c.logger.Debug().Str("filter", filter).Int("num_timeseries", len(timeSeriesList)).Msg("processing time series list")
	for _, timeSeries := range timeSeriesList {
		var tags circonus.Tags
		tags = append(tags, baseTags...)
		tags = append(tags, circonus.Tag{Category: "resource_type", Value: timeSeries.GetResource().GetType()})
		for k, v := range timeSeries.GetResource().GetLabels() {
			// c.logger.Debug().Str("category", k).Msg("adding resource label")
			tags = append(tags, circonus.Tag{Category: k, Value: v})
		}
		for k, v := range timeSeries.GetMetric().GetLabels() {
			if k == "metric_type" {
				continue
			}
			if k == "metric_kind" {
				continue
			}
			// c.logger.Debug().Str("category", k).Msg("adding metric label")
			tags = append(tags, circonus.Tag{Category: k, Value: v})
		}

		metricType := ""
		switch timeSeries.GetValueType().String() {
		case "DOUBLE":
			metricType = circonus.MetricTypeFloat64
		case "INT64":
			metricType = circonus.MetricTypeUint64
		case "STRING":
			metricType = circonus.MetricTypeString
		default:
			c.logger.Warn().Str("type", timeSeries.GetValueType().String()).Msg("unmapped metric value type, ignoring")
			continue
		}

		mn := c.check.MetricNameWithStreamTags(metricName, tags)

		// NOTE: The broker seems to want repeated metrics in oldest to newest order.
		//       The GCP TimeSeriesList API no longer supports requesting in a specific
		//       order - the points are returned in (reverse) newest to oldest order.
		// for _, pt := range timeSeries.Points {
		pts := timeSeries.GetPoints()
		for i := len(pts) - 1; i >= 0; i-- {
			pt := pts[i]
			if c.done() {
				break
			}
			var value interface{}
			end := pt.GetInterval().GetEndTime()
			ts := time.Unix(end.GetSeconds(), int64(end.GetNanos()))
			switch metricType {
			case circonus.MetricTypeFloat64:
				value = pt.GetValue().GetDoubleValue()
			case circonus.MetricTypeString:
				value = pt.GetValue().GetStringValue()
			case circonus.MetricTypeUint64:
				value = pt.GetValue().GetInt64Value()
			default:
				c.logger.Error().Str("name", metricName).Str("type", metricType).Msg("invalid metric type")
				continue
			}
			if err := c.check.WriteMetricSample(metricDest, mn, metricType, value, &ts); err != nil {
				c.logger.Warn().Err(err).Str("name", metricName).Msg("recording metric")
			}
		}
		if c.done() {
			break
		}
	}
}
