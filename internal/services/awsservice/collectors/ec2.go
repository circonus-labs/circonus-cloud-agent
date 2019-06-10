// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"bytes"
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// handle AWS/EC2 specific tasks
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-cloudwatch.html

// EC2 defines the collector instance
type EC2 struct {
	common
	filters *[]Filter
}

// ec2instance is an internal structure which contains the aws InstanceId
// and a list of base stream tags from the ec2 instance's meta data
type ec2instance struct {
	id   string
	tags circonus.Tags
}

func newEC2(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/EC2"
	c := &EC2{
		common:  newCommon(ctx, ns, check, cfg, logger),
		filters: cfg.InstanceFilters,
	}
	if len(c.metrics) == 0 {
		c.metrics = c.DefaultMetrics()
	}
	c.tags = append(c.tags, circonus.Tag{Category: "service", Value: ns})
	c.logger.Debug().Msg("initialized")
	return c, nil
}

// Collect pulls list of active ec2 instances, then configured metrics from
// cloudwatch, forwarding them to circonus.
func (c *EC2) Collect(sess *session.Session, timespan MetricTimespan, baseTags circonus.Tags) error {
	if sess == nil {
		return errors.New("invalid session (nil)")
	}

	if !c.Enabled() {
		return nil
	}

	c.logger.Debug().Msg("getting aws ec2 instance list")
	ec2instances, err := c.ec2Instances(sess, baseTags)
	if awserr := c.trackAWSErrors(err); awserr != nil {
		return errors.Wrap(c.trackAWSErrors(awserr), "geting instance information")
	}

	c.logger.Debug().Msg("retrieving telemetry")

	// NOTE: remove need for custom metricStats and metricData collectors
	// collectorFn := c.ec2MetricStats
	// if c.useGMD {
	// 	collectorFn = c.ec2MetricData
	// }
	collectorFn := c.metricStats
	if c.useGMD {
		collectorFn = c.metricData
	}
	var buf bytes.Buffer
	buf.Grow(32768)
	metricDimensionName := "InstanceId"
	for _, instanceInfo := range ec2instances {
		dims := []*cloudwatch.Dimension{
			{
				Name:  &metricDimensionName,
				Value: &instanceInfo.id,
			},
		}
		var metricTags []circonus.Tag
		if len(baseTags) > 0 {
			metricTags = append(metricTags, baseTags...)
		}
		if len(instanceInfo.tags) > 0 {
			metricTags = append(metricTags, instanceInfo.tags...)
		}
		if err := collectorFn(&buf, sess, timespan, dims, metricTags); err != nil {
			c.logger.Error().Err(err).Msg("collecting telemetry")
		}
		if buf.Len() == 0 {
			c.logger.Warn().Str("collector", c.ID()).Msg("no telemetry to submit")
			continue
		}
		c.logger.Debug().Str("collector", c.ID()).Msg("submitting telemetry")
		if err := c.check.SubmitMetrics(&buf); err != nil {
			c.logger.Error().Err(err).Msg("submitting telemetry")
		}
		buf.Reset()
	}

	return nil
}

// ec2Instances pulls a list of ec2 instances, saves the InstanceId for the
// cloudwatch metric dimension and creates a list of default stream tags to
// use for the metrics collected for the specific ec2 instance.
func (c *EC2) ec2Instances(sess *session.Session, baseTags circonus.Tags) ([]ec2instance, error) {
	ec2List := []ec2instance{}

	if sess == nil {
		return ec2List, errors.New("invalid session (nil)")
	}

	ec2Svc := ec2.New(sess)
	var describeInstancesInput *ec2.DescribeInstancesInput
	if c.filters != nil && len(*c.filters) > 0 {
		filters := make([]*ec2.Filter, len(*c.filters))
		for idx, filter := range *c.filters {
			filters[idx] = &ec2.Filter{Name: filter.Name, Values: filter.Values}
		}
		describeInstancesInput = &ec2.DescribeInstancesInput{Filters: filters}
	}
	results, err := ec2Svc.DescribeInstances(describeInstancesInput)
	if err != nil {
		return ec2List, errors.Wrap(err, "describing instances")
	}

	if len(c.tags) > 0 {
		baseTags = append(baseTags, c.tags...)
	}
	for _, reservation := range results.Reservations {
		for _, ec2inst := range reservation.Instances {
			if *ec2inst.State.Name != "running" {
				continue
			}

			streamTags := circonus.Tags{
				circonus.Tag{Category: "zone", Value: *ec2inst.Placement.AvailabilityZone},
			}
			streamTags = append(streamTags, baseTags...)
			streamTags = append(streamTags, circonus.Tags{
				circonus.Tag{Category: "type", Value: *ec2inst.InstanceType},
				circonus.Tag{Category: "arch", Value: *ec2inst.Architecture},
				circonus.Tag{Category: "image_id", Value: *ec2inst.ImageId},
			}...)
			if len((*ec2inst).Tags) > 0 {
				for _, tag := range (*ec2inst).Tags {
					tc := strings.ToLower(strings.Replace(*tag.Key, ":", "_", -1))
					tv := strings.ToLower(*tag.Value)
					streamTags = append(streamTags, circonus.Tag{Category: tc, Value: tv})
				}
			}
			ec2List = append(ec2List, ec2instance{id: *ec2inst.InstanceId, tags: streamTags})
			if c.done() {
				return []ec2instance{}, nil
			}
		}
	}

	return ec2List, nil
}

// // nolint: gocyclo
// // ec2MetricData uses cloudwatch.GetMetricData (one call per 100 metrics)
// func (c *EC2) ec2MetricData(metricDest io.Writer, sess *session.Session, timespan MetricTimespan, dimensions []*cloudwatch.Dimension, baseTags circonus.Tags) error {
// 	if metricDest == nil {
// 		return errors.New("invalid metric destination (nil)")
// 	}
// 	if sess == nil {
// 		return errors.New("invalid session (nil)")
// 	}

// 	metricDataQueryBuckets := [][]*cloudwatch.MetricDataQuery{make([]*cloudwatch.MetricDataQuery, 0, 100)}
// 	bucketID := 0
// 	metricCount := 0
// 	returnData := true

// 	for metricIdx, metricDefinition := range c.metrics {
// 		if metricDefinition.AWSMetric.Disabled {
// 			continue
// 		}
// 		for _, metricStatName := range metricDefinition.AWSMetric.Stats {
// 			metricStatName := metricStatName
// 			metricID := fmt.Sprintf(resultIDFormat, metricIdx, metricStatName, metricCount)
// 			metricStat := cloudwatch.MetricStat{
// 				Metric: &cloudwatch.Metric{
// 					MetricName: &metricDefinition.AWSMetric.Name,
// 					Namespace:  &c.id,
// 				},
// 				Period: &timespan.Period,
// 				Stat:   &metricStatName,
// 			}
// 			metricStat.Metric.Dimensions = append(metricStat.Metric.Dimensions, dimensions...)
// 			if len(c.dimensions) > 0 {
// 				metricStat.Metric.Dimensions = append(metricStat.Metric.Dimensions, c.dimensions...)
// 			}
// 			metricDataQuery := &cloudwatch.MetricDataQuery{
// 				Id:         &metricID,
// 				ReturnData: &returnData,
// 				MetricStat: &metricStat,
// 			}

// 			metricDataQueryBuckets[bucketID] = append(metricDataQueryBuckets[bucketID], metricDataQuery)
// 			metricCount++

// 			// metric data query capped at 100 metrics, create new bucket every 100 metrics
// 			if metricCount%100 == 0 {
// 				bucketID++
// 				metricDataQueryBuckets = append(metricDataQueryBuckets, make([]*cloudwatch.MetricDataQuery, 0, 100))
// 			}
// 			if c.done() {
// 				return nil
// 			}
// 		}
// 	}

// 	cwSvc := cloudwatch.New(sess)

// 	for _, metricDataQueries := range metricDataQueryBuckets {
// 		getMetricDataInput := cloudwatch.GetMetricDataInput{
// 			StartTime:         &timespan.Start,
// 			EndTime:           &timespan.End,
// 			MetricDataQueries: metricDataQueries,
// 		}
// 		results, err := cwSvc.GetMetricData(&getMetricDataInput)
// 		if err != nil {
// 			c.logger.Error().Err(err).Msg("retrieving metric data")
// 			continue
// 		}
// 		if c.done() {
// 			return nil
// 		}
// 		for {
// 			for _, result := range results.MetricDataResults {
// 				var metricStat string
// 				var metricIdx, queryIdx int
// 				if n, err := fmt.Sscanf(*result.Id, resultIDFormat, &metricIdx, &metricStat, &queryIdx); err != nil {
// 					c.logger.Error().Err(err).Str("result_id", *result.Id).Msg("unable to extract cance/metric IDs from result id")
// 					continue
// 				} else if n != 2 {
// 					c.logger.Error().Int("num_extracted", n).Str("result_id", *result.Id).Msg("unable to extract BOTH instance id and metric id from result id")
// 					continue
// 				}

// 				if metricIdx > len(c.metrics) || metricIdx < 0 {
// 					c.logger.Error().Int("metric_idx", metricIdx).Int("num_metrics", len(c.metrics)).Msg("invalid metric index <0||>len")
// 					continue
// 				}

// 				if metricStat == "" {
// 					c.logger.Error().Str("result_id", *result.Id).Msg("invalid metric stat")
// 					continue
// 				}

// 				if queryIdx > len(metricDataQueries) || queryIdx < 0 {
// 					c.logger.Error().Int("query_idx", queryIdx).Int("num_queries", len(metricDataQueries)).Msg("invalid metric data query index <0||>len")
// 					continue
// 				}

// 				metricDefinition := c.metrics[metricIdx]

// 				var metricTags []circonus.Tag
// 				if len(baseTags) > 0 {
// 					metricTags = append(metricTags, baseTags...)
// 				}
// 				if len(metricDataQueries[queryIdx].MetricStat.Metric.Dimensions) > 0 {
// 					for _, d := range metricDataQueries[queryIdx].MetricStat.Metric.Dimensions {
// 						metricTags = append(metricTags, circonus.Tag{Category: *d.Name, Value: *d.Value})
// 					}
// 				}

// 				for idx, resultTimestamp := range result.Timestamps {
// 					metricValue := *result.Values[idx]
// 					if err := c.recordMetric(metricDest, metricDefinition, metricStat, metricValue, resultTimestamp, metricTags); err != nil {
// 						c.logger.Warn().Err(err).Str("aws_metric", metricDefinition.AWSMetric.Name).Msg("recording metric")
// 					}
// 				}
// 				if c.done() {
// 					return nil
// 				}
// 			}

// 			if results.NextToken == nil {
// 				break
// 			}
// 			getMetricDataInput.SetNextToken(*results.NextToken)
// 			results, err = cwSvc.GetMetricData(&getMetricDataInput)
// 			if err != nil {
// 				c.logger.Error().Err(err).Msg("retrieving metric data w/NextToken")
// 				break
// 			}
// 		}
// 	}

// 	return nil
// }

// func (c *EC2) ec2MetricStats(metricDest io.Writer, sess *session.Session, timespan MetricTimespan, dimensions []*cloudwatch.Dimension, baseTags circonus.Tags) error {
// 	if metricDest == nil {
// 		return errors.New("invalid metric destination (nil)")
// 	}
// 	if sess == nil {
// 		return errors.New("invalid session (nil)")
// 	}

// 	cwSvc := cloudwatch.New(sess)
// 	for _, metricDefinition := range c.metrics {
// 		if metricDefinition.AWSMetric.Disabled {
// 			continue
// 		}
// 		stats := make([]*string, len(metricDefinition.AWSMetric.Stats))
// 		for i := range metricDefinition.AWSMetric.Stats {
// 			stats[i] = &metricDefinition.AWSMetric.Stats[i]
// 		}
// 		getMetricStatisticsInput := cloudwatch.GetMetricStatisticsInput{
// 			Statistics: stats,
// 			Dimensions: dimensions,
// 			MetricName: &metricDefinition.AWSMetric.Name,
// 			Namespace:  &c.id,
// 			StartTime:  &timespan.Start,
// 			EndTime:    &timespan.End,
// 			Period:     &timespan.Period,
// 		}
// 		if len(c.dimensions) > 0 {
// 			getMetricStatisticsInput.Dimensions = append(getMetricStatisticsInput.Dimensions, c.dimensions...)
// 		}

// 		result, err := cwSvc.GetMetricStatistics(&getMetricStatisticsInput)
// 		if err != nil {
// 			c.logger.Error().Err(err).Msg("retrieving metric data")
// 			continue
// 		}

// 		var metricTags circonus.Tags
// 		if len(baseTags) > 0 {
// 			metricTags = append(metricTags, baseTags...)
// 		}
// 		if len(getMetricStatisticsInput.Dimensions) > 0 {
// 			for _, d := range getMetricStatisticsInput.Dimensions {
// 				metricTags = append(metricTags, circonus.Tag{Category: *d.Name, Value: *d.Value})
// 			}
// 		}

// 		for _, metricDatapoint := range result.Datapoints {
// 			for _, stat := range metricDefinition.AWSMetric.Stats {
// 				var metricValue float64
// 				switch stat {
// 				case metricStatAverage:
// 					metricValue = *metricDatapoint.Average
// 				case metricStatSum:
// 					metricValue = *metricDatapoint.Sum
// 				case metricStatMinimum:
// 					metricValue = *metricDatapoint.Minimum
// 				case metricStatMaximum:
// 					metricValue = *metricDatapoint.Maximum
// 				case metricStatSampleCount:
// 					metricValue = *metricDatapoint.SampleCount
// 				}
// 				if err := c.recordMetric(metricDest, metricDefinition, stat, metricValue, metricDatapoint.Timestamp, metricTags); err != nil {
// 					c.logger.Warn().Err(err).Str("aws_metric", metricDefinition.AWSMetric.Name).Msg("recording metric")
// 				}
// 				if c.done() {
// 					return nil
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

// DefaultMetrics defines the default EC2 metrics
func (c *EC2) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric{
				Name:  "CPUUtilization",
				Stats: []string{metricStatAverage},
				Units: "Percent",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "DiskReadOps",
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
				Name:  "DiskWriteOps",
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
				Name:  "DiskReadBytes",
				Stats: []string{metricStatAverage},
				Units: "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "DiskWriteBytes",
				Stats: []string{metricStatAverage},
				Units: "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically

			},
		},
		{
			AWSMetric{
				Name:  "NetworkIn",
				Stats: []string{metricStatAverage},
				Units: "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "NetworkOut",
				Stats: []string{metricStatAverage},
				Units: "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "NetworkPacketsIn",
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
				Name:  "NetworkPacketsOut",
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
				Name:  "EBSReadOps",
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
				Name:  "EBSWriteOps",
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
				Name:  "EBSReadBytes",
				Stats: []string{metricStatAverage},
				Units: "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Name:  "EBSWriteBytes",
				Stats: []string{metricStatAverage},
				Units: "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
	}
}
