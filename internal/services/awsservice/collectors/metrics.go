// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"fmt"
	"io"
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
				if n, err := fmt.Sscanf(*result.Id, resultIDFormat, &metricIdx, &metricStat, &queryIdx); err != nil {
					c.logger.Error().Err(err).Str("result_id", *result.Id).Msg("unable to extract cance/metric IDs from result id")
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
				for idx, resultTimestamp := range result.Timestamps {
					metricValue := *result.Values[idx]
					if err := c.recordMetric(metricDest, metricDefinition, metricStat, metricValue, resultTimestamp, metricTags); err != nil {
						c.logger.Warn().Err(err).Str("aws_metric", metricDefinition.AWSMetric.Name).Msg("recording metric datapoint")
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
			getMetricStatisticsInput.Dimensions = c.dimensions
		} else if len(c.dimensions) > 0 {
			getMetricStatisticsInput.Dimensions = c.dimensions
		}

		result, err := cwSvc.GetMetricStatistics(&getMetricStatisticsInput)
		if err != nil {
			c.logger.Error().Err(err).Str("aws_metric_name", metricDefinition.AWSMetric.Name).Msg("retrieving metric statistics")
			continue
		}
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
		for _, metricDatapoint := range result.Datapoints {
			for _, stat := range metricDefinition.AWSMetric.Stats {
				var metricValue float64
				switch stat {
				case "Average":
					metricValue = *metricDatapoint.Average
				case "Sum":
					metricValue = *metricDatapoint.Sum
				case "Minimum":
					metricValue = *metricDatapoint.Minimum
				case "Maximum":
					metricValue = *metricDatapoint.Maximum
				case "SampleCount":
					metricValue = *metricDatapoint.SampleCount
				}
				if err := c.recordMetric(metricDest, metricDefinition, stat, metricValue, metricDatapoint.Timestamp, metricTags); err != nil {
					c.logger.Warn().Err(err).Str("aws_metric", metricDefinition.AWSMetric.Name).Msg("recording metric statistic")
				}
				if c.done() {
					return nil
				}
			}
		}
	}

	return nil
}

// recordMetric creates a metric name w/encoded stream tags then writes the metric sample to the metric destination
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
		err = c.check.WriteMetricSample(metricDest, metricName, "n", val.(float64), ts)
	case "histogram":
		err = c.check.WriteMetricSample(metricDest, metricName, "n", val.(float64), ts)
	case "text":
		err = c.check.WriteMetricSample(metricDest, metricName, "s", fmt.Sprintf("%v", val), ts)
	default:
		c.logger.Warn().Interface("metric", metric).Msg("invalid Circonus Metric Type configured, ignoring metric sample")
	}

	return err
}
