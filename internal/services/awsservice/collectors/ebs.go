// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"bytes"
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// handle AWS/EBS specific tasks
// updated doc link: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using_cloudwatch_ebs.html
// original doc link (no longer works): https://docs.aws.amazon.com/elasticloadbalancing/latest/classic/EBS-cloudwatch-metrics.html

// EBS defines the collector instance
type EBS struct {
	common
}

func newEBS(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/EBS"
	c := &EBS{
		common: newCommon(ctx, ns, check, cfg, logger),
	}
	if len(c.metrics) == 0 {
		c.metrics = c.DefaultMetrics()
	}
	c.tags = append(c.tags, circonus.Tag{Category: "service", Value: ns})
	c.logger.Debug().Msg("initialized")
	return c, nil
}

// Collect pulls list of ebs volumes from ec2 instances, then configured metrics from
// cloudwatch, forwarding them to circonus.
func (c *EBS) Collect(sess *session.Session, timespan MetricTimespan, baseTags circonus.Tags) error {
	if sess == nil {
		return errors.New("invalid session (nil)")
	}

	if !c.Enabled() {
		return nil
	}

	c.logger.Debug().Msg("retrieving telemetry")

	collectorFn := c.metricStats
	if c.useGMD {
		collectorFn = c.metricData
	}

	var buf bytes.Buffer
	buf.Grow(32768)

	// call once for entire zone
	var dims []*cloudwatch.Dimension
	if err := collectorFn(&buf, sess, timespan, dims, baseTags); err != nil {
		return errors.Wrap(err, "collecting telemetry")
	}
	if buf.Len() > 0 {
		c.logger.Debug().Str("collector", c.ID()).Msg("submitting telemetry")
		if err := c.check.SubmitMetrics(&buf); err != nil {
			c.logger.Error().Err(err).Msg("submitting telemetry")
		}
		buf.Reset()
	}

	// call for each volume returned from query to ec2 service
	c.logger.Debug().Msg("getting aws ec2 instance ebs volume list")
	ebsVolumes, err := c.ebsVolumes(sess, baseTags)
	if awserr := c.trackAWSErrors(err); awserr != nil {
		return errors.Wrap(c.trackAWSErrors(awserr), "getting ebs volume information")
	}
	if len(ebsVolumes) == 0 {
		c.logger.Warn().Msg("zero ebs volumes found")
	}
	metricDimensionName := "VolumeId"
	for _, volumeInfo := range ebsVolumes {
		dims := []*cloudwatch.Dimension{
			{
				Name:  &metricDimensionName,
				Value: &volumeInfo.id,
			},
		}
		var metricTags []circonus.Tag
		if len(baseTags) > 0 {
			metricTags = append(metricTags, baseTags...)
		}
		if len(volumeInfo.tags) > 0 {
			metricTags = append(metricTags, volumeInfo.tags...)
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

type ebsVolume struct {
	id   string
	tags circonus.Tags
}

// ebsVolumes lists ec2 instance volumes, saves the VolumeId for the
// cloudwatch metric dimension and creates a list of default stream tags to
// use for the metrics collected for the specific ebs volume.
func (c *EBS) ebsVolumes(sess client.ConfigProvider, baseTags circonus.Tags) ([]ebsVolume, error) {
	ebsVolumes := []ebsVolume{}

	if sess == nil {
		return ebsVolumes, errors.New("invalid session (nil)")
	}

	ec2Svc := ec2.New(sess)
	input := &ec2.DescribeVolumesInput{}
	results, err := ec2Svc.DescribeVolumes(input)
	if err != nil {
		return ebsVolumes, errors.Wrap(err, "describing ebs volumes")
	}

	if len(c.tags) > 0 {
		baseTags = append(baseTags, c.tags...)
	}
	for _, volume := range results.Volumes {
		vid := *volume.VolumeId
		streamTags := circonus.Tags{
			circonus.Tag{Category: "type", Value: *volume.VolumeType},
			circonus.Tag{Category: "zone", Value: *volume.AvailabilityZone},
		}
		streamTags = append(streamTags, baseTags...)
		if len(volume.Tags) > 0 {
			for _, tag := range volume.Tags {
				cat := strings.ToLower(strings.ReplaceAll(*tag.Key, ":", "_"))
				val := strings.ToLower(*tag.Value)
				streamTags = append(streamTags, circonus.Tag{Category: cat, Value: val})
			}
		}
		ebsVolumes = append(ebsVolumes, ebsVolume{id: vid, tags: streamTags})
		if c.done() {
			return []ebsVolume{}, nil
		}
	}

	return ebsVolumes, nil
}

// DefaultMetrics returns a default metric configuration
func (c *EBS) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric{
				Name:  "VolumeReadBytes",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
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
				Name:  "VolumeWriteBytes",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
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
				Name:  "VolumeReadOps",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
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
				Name:  "VolumeWriteOps",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
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
				Name:  "VolumeTotalReadTime",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
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
				Name:  "VolumeTotalWriteTime",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
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
				Name:  "VolumeIdleTime",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
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
				Name:  "VolumeQueueLength",
				Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
				Units: "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		//
		// Use with provisioned iops ssd volumns only
		//
		// {
		// 	AWSMetric{
		// 		Name:  "VolumeThroughputPercentage",
		// 		Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
		// 		Units: "Percent",
		// 	},
		// 	CirconusMetric{
		// 		Name: "",              // NOTE: AWSMetric.Name will be used if blank
		// 		Type: "gauge",         // (gauge|counter|histogram|text)
		// 		Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
		// 	},
		// },
		// {
		// 	AWSMetric{
		// 		Name:  "VolumeConsumedReadWriteOps",
		// 		Stats: []string{metricStatMinimum, metricStatMaximum, metricStatAverage, metricStatSum, metricStatSampleCount},
		// 		Units: metricStatSum,
		// 	},
		// 	CirconusMetric{
		// 		Name: "",              // NOTE: AWSMetric.Name will be used if blank
		// 		Type: "gauge",         // (gauge|counter|histogram|text)
		// 		Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
		// 	},
		// },
		//
		// Use with general purpose ssd (gp2), throughput optimized hdd (st1), and cold hdd (sc1) volumes only
		//
		// {
		// 	AWSMetric{
		// 		Name:  "BurstBalance",
		// 		Stats: []string{metricStatAverage},
		// 		Units: "Percent",
		// 	},
		// 	CirconusMetric{
		// 		Name: "",              // NOTE: AWSMetric.Name will be used if blank
		// 		Type: "gauge",         // (gauge|counter|histogram|text)
		// 		Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
		// 	},
		// },
	}
}
