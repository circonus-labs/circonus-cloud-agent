// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/pkg/errors"
)

// Generic metric collection methods
//   NOTE: override in specific service if needed (see ec2)

// nolint: gocyclo
func (c *common) metricData(metricDest io.Writer, sess client.ConfigProvider, timespan MetricTimespan, dimensions []*cloudwatch.Dimension, baseTags circonus.Tags) error {
	if metricDest == nil {
		return errors.New("invalid metric destination (nil)")
	}
	if sess == nil {
		return errors.New("invalid session (nil)")
	}

	metricDataQueryBuckets := [][]*cloudwatch.MetricDataQuery{make([]*cloudwatch.MetricDataQuery, 0, 100)}
	bucketID := 0
	metricCount := 0
	returnData := true
	for metricIdx, metricDefinition := range c.metrics {
		if metricDefinition.AWSMetric.Disabled {
			continue
		}

		metricDefinition := metricDefinition

		for _, metricStatName := range metricDefinition.AWSMetric.Stats {
			metricStatName := metricStatName
			metricID := fmt.Sprintf(resultIDFormat, metricIdx, metricStatName, metricCount)
			metricStat := cloudwatch.MetricStat{
				Metric: &cloudwatch.Metric{
					MetricName: &metricDefinition.AWSMetric.Name,
					Namespace:  &c.id,
				},
				Period: &timespan.Period,
				Stat:   &metricStatName,
			}
			if len(dimensions) > 0 {
				metricStat.Metric.Dimensions = c.dimensions
			} else if len(c.dimensions) > 0 {
				metricStat.Metric.Dimensions = c.dimensions
			}

			metricDataQuery := &cloudwatch.MetricDataQuery{
				Id:         &metricID,
				ReturnData: &returnData,
				MetricStat: &metricStat,
			}

			metricDataQueryBuckets[bucketID] = append(metricDataQueryBuckets[bucketID], metricDataQuery)
			metricCount++

			// metric data query capped at 100 metrics, create new bucket every 100 metrics
			if metricCount%100 == 0 {
				bucketID++
				metricDataQueryBuckets = append(metricDataQueryBuckets, make([]*cloudwatch.MetricDataQuery, 0, 100))
			}
			if c.done() {
				return nil
			}
		}
	}

	cwSvc := cloudwatch.New(sess)

	for _, metricDataQueries := range metricDataQueryBuckets {
		getMetricDataInput := cloudwatch.GetMetricDataInput{
			StartTime:         &timespan.Start,
			EndTime:           &timespan.End,
			MetricDataQueries: metricDataQueries,
		}
		results, err := cwSvc.GetMetricData(&getMetricDataInput)
		if err != nil {
			c.logger.Error().Err(err).Msg("retrieving metric data")
			continue
		}
		if c.done() {
			return nil
		}
		for {
			for _, result := range results.MetricDataResults {
				var metricStat string
				var metricIdx, queryIdx int
				if n, err2 := fmt.Sscanf(*result.Id, resultIDFormat, &metricIdx, &metricStat, &queryIdx); err2 != nil {
					c.logger.Error().Err(err2).Str("result_id", *result.Id).Msg("unable to extract metric IDs from result id")
					continue
				} else if n != 2 {
					c.logger.Error().Int("num_extracted", n).Str("result_id", *result.Id).Msg("unable to extract BOTH instance id and metric id from result id")
					continue
				}

				if metricIdx > len(c.metrics) || metricIdx < 0 {
					c.logger.Error().Int("metric_idx", metricIdx).Int("num_metrics", len(c.metrics)).Msg("invalid metric index <0||>len")
					continue
				}

				if metricStat == "" {
					c.logger.Error().Str("result_id", *result.Id).Msg("invalid metric stat")
					continue
				}

				if queryIdx > len(metricDataQueries) || queryIdx < 0 {
					c.logger.Error().Int("query_idx", queryIdx).Int("num_queries", len(metricDataQueries)).Msg("invalid metric data query index <0||>len")
					continue
				}

				var metricTags circonus.Tags
				if len(c.tags) > 0 {
					metricTags = append(metricTags, c.tags...)
				}
				if len(baseTags) > 0 {
					metricTags = append(metricTags, baseTags...)
				}
				if len(metricDataQueries[queryIdx].MetricStat.Metric.Dimensions) > 0 {
					for _, d := range metricDataQueries[queryIdx].MetricStat.Metric.Dimensions {
						metricTags = append(metricTags, circonus.Tag{Category: *d.Name, Value: *d.Value})
					}
				}

				metricDefinition := c.metrics[metricIdx]
				samples := c.sortMetricDataSamples(result)
				for _, sample := range samples {
					if err2 := c.recordMetric(metricDest, metricDefinition, metricStat, sample.Value, sample.TS, metricTags); err2 != nil {
						c.logger.Warn().Err(err2).Str("aws_metric", metricDefinition.AWSMetric.Name).Msg("recording metric data point")
					}
				}

				if c.done() {
					return nil
				}
			}

			if results.NextToken == nil {
				break
			}
			getMetricDataInput.SetNextToken(*results.NextToken)
			results, err = cwSvc.GetMetricData(&getMetricDataInput)
			if err != nil {
				c.logger.Error().Err(err).Msg("retrieving metric data w/NextToken")
				break
			}
		}
	}

	return nil
}

// The following is for sorting the time series samples
// returned by aws GetMetricData

type mdSample struct {
	TS    *time.Time
	Value float64
}

func (c *common) sortMetricDataSamples(ts *cloudwatch.MetricDataResult) []mdSample {
	if len(ts.Timestamps) == 0 {
		return []mdSample{}
	}

	samples := make([]mdSample, len(ts.Timestamps))
	for idx, rts := range ts.Timestamps {
		samples[idx] = mdSample{TS: rts, Value: *ts.Values[idx]}
	}

	sort.Slice(samples, func(i, j int) bool { return samples[i].TS.Before(*samples[j].TS) })
	return samples
}

func (c *common) metricStats(metricDest io.Writer, sess client.ConfigProvider, timespan MetricTimespan, dimensions []*cloudwatch.Dimension, baseTags circonus.Tags) error {
	if metricDest == nil {
		return errors.New("invalid metric destination (nil)")
	}
	if sess == nil {
		return errors.New("invalid session (nil)")
	}

	cwSvc := cloudwatch.New(sess)

	for _, metricDefinition := range c.metrics {
		if metricDefinition.AWSMetric.Disabled {
			continue
		}

		metricDefinition := metricDefinition

		stats := make([]*string, len(metricDefinition.AWSMetric.Stats))
		for i := range metricDefinition.AWSMetric.Stats {
			stats[i] = &metricDefinition.AWSMetric.Stats[i]
		}
		getMetricStatisticsInput := cloudwatch.GetMetricStatisticsInput{
			Statistics: stats,
			MetricName: &metricDefinition.AWSMetric.Name,
			Namespace:  &c.id,
			StartTime:  &timespan.Start,
			EndTime:    &timespan.End,
			Period:     &timespan.Period,
		}
		if len(dimensions) > 0 {
			getMetricStatisticsInput.Dimensions = dimensions
		} else if len(c.dimensions) > 0 {
			getMetricStatisticsInput.Dimensions = c.dimensions
		}

		c.logger.Debug().Interface("inputs", getMetricStatisticsInput).Msg("metric stats inputs")

		result, err := cwSvc.GetMetricStatistics(&getMetricStatisticsInput)
		if err != nil {
			c.logger.Error().Err(err).Str("aws_metric_name", metricDefinition.AWSMetric.Name).Msg("retrieving metric statistics")
			continue
		}
		c.logger.Debug().Interface("result", result).Str("metric", metricDefinition.AWSMetric.Name).Msg("AWS response")
		var metricTags circonus.Tags
		if len(c.tags) > 0 {
			metricTags = append(metricTags, c.tags...)
		}
		if len(baseTags) > 0 {
			metricTags = append(metricTags, baseTags...)
		}
		if len(getMetricStatisticsInput.Dimensions) > 0 {
			for _, d := range getMetricStatisticsInput.Dimensions {
				metricTags = append(metricTags, circonus.Tag{Category: *d.Name, Value: *d.Value})
			}
		}
		datapoints := c.sortMetricStatDatapoints(result.Datapoints, metricDefinition)
		for _, dp := range datapoints {
			var mt circonus.Tags
			mt = append(mt, metricTags...)
			mt = append(mt, circonus.Tag{Category: "units", Value: dp.Units})
			if err := c.recordMetric(metricDest, metricDefinition, dp.Stat, dp.Value, dp.Timestamp, mt); err != nil {
				c.logger.Warn().Err(err).Str("aws_metric", metricDefinition.AWSMetric.Name).Msg("recording metric statistic")
			}
		}
		if c.done() {
			return nil
		}
	}

	return nil
}

// The following is for sorting the datapoints (in time stamp order)
// returned by aws GetMetricStatistics

type msDatapoint struct {
	Timestamp *time.Time
	Units     string
	Stat      string
	Value     float64
}

func (c *common) sortMetricStatDatapoints(datapoints []*cloudwatch.Datapoint, md Metric) []msDatapoint {
	if len(datapoints) == 0 {
		return []msDatapoint{}
	}

	samples := make([]msDatapoint, 0)
	for _, dp := range datapoints {
		for _, stat := range md.AWSMetric.Stats {
			var v float64
			switch stat {
			case "Average":
				v = *dp.Average
			case "Sum":
				v = *dp.Sum
			case "Minimum":
				v = *dp.Minimum
			case "Maximum":
				v = *dp.Maximum
			case "SampleCount":
				v = *dp.SampleCount
			}

			samples = append(samples, msDatapoint{
				Timestamp: dp.Timestamp,
				Value:     v,
				Units:     *dp.Unit,
				Stat:      stat,
			})
			if c.done() {
				break
			}
		}
		if c.done() {
			break
		}
	}

	sort.Slice(samples, func(i, j int) bool { return samples[i].Timestamp.Before(*samples[j].Timestamp) })
	return samples
}

// recordMetric creates a metric name w/encoded stream tags then writes the metric sample to the metric destination.
func (c *common) recordMetric(metricDest io.Writer, metric Metric, metricStat string, val interface{}, ts *time.Time, baseTags circonus.Tags) error {
	mn := metric.CirconusMetric.Name
	if mn == "" {
		mn = metric.AWSMetric.Name
	}
	if metricStat != "" {
		mn += "`" + metricStat
	}

	var tags circonus.Tags
	tags = append(tags, baseTags...)

	if len(metric.CirconusMetric.Tags) > 0 {
		tags = append(tags, metric.CirconusMetric.Tags...)
	}
	if metric.AWSMetric.Units != "" {
		tags = append(tags, circonus.Tag{Category: "units", Value: strings.ToLower(metric.AWSMetric.Units)})
	}

	metricName := c.check.MetricNameWithStreamTags(mn, tags)

	var err error
	switch metric.CirconusMetric.Type {
	case "counter":
		fallthrough
	case "gauge":
		if strings.Contains(metricName, "CPUUtilization") {
			mt := c.check.EncodeMetricTags(tags)
			c.logger.Debug().Str("encoded_metric_name", metricName).Int64("epoch", ts.Unix()).Msg("for data api call")
			c.logger.Debug().Str("metric", mn).Strs("tags", mt).Str("type", "n").Float64("val", val.(float64)).Time("ts", *ts).Msg("metric to circonus")
		}
		err = c.check.WriteMetricSample(metricDest, metricName, "n", val.(float64), ts)
	case "histogram":
		if strings.Contains(metricName, "CPUUtilization") {
			mt := c.check.EncodeMetricTags(tags)
			c.logger.Debug().Str("encoded_metric_name", metricName).Int64("epoch", ts.Unix()).Msg("for data api call")
			c.logger.Debug().Str("metric", mn).Strs("tags", mt).Str("type", "h").Float64("val", val.(float64)).Time("ts", *ts).Msg("metric to circonus")
		}
		err = c.check.WriteMetricSample(metricDest, metricName, "h", val.(float64), ts)
	case "text":
		err = c.check.WriteMetricSample(metricDest, metricName, "s", fmt.Sprintf("%v", val), ts)
	default:
		c.logger.Warn().Interface("metric", metric).Msg("invalid Circonus Metric Type configured, ignoring metric sample")
	}

	return err
}
