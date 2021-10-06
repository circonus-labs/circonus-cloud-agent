// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// handle AWS/ElastiCache specific tasks
// https://docs.aws.amazon.com/AWSElastiCache/latest/UserGuide/using-cloudwatch.html

// ElastiCache defines the collector instance.
type ElastiCache struct {
	clusterIDs *[]string
	common
}

// newElastiCache creates a new ElastiCache telemetry collector.
func newElastiCache(ctx context.Context, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) (Collector, error) {
	ns := "AWS/ElastiCache"
	c := &ElastiCache{
		common:     newCommon(ctx, ns, check, cfg, logger),
		clusterIDs: cfg.CacheClusterIDs,
	}
	if len(c.metrics) == 0 {
		c.metrics = c.DefaultMetrics()
	}
	c.tags = append(c.tags, circonus.Tag{Category: "service", Value: ns})
	c.logger.Debug().Msg("initialized")
	return c, nil
}

// Collect pulls list of active clusters and nodes for each cluster (the
// dimensions required to collect telemetry), then configured metrics from
// cloudwatch, then collects each enable telmetry point forwarding them
// to circonus.
func (c *ElastiCache) Collect(sess *session.Session, timespan MetricTimespan, baseTags circonus.Tags) error {
	if sess == nil {
		return errors.New("invalid session (nil)")
	}

	if !c.Enabled() {
		return nil
	}

	c.logger.Debug().Msg("retrieving telemetry")
	clusterList, err := c.clusterList(sess)
	if awserr := c.trackAWSErrors(err); awserr != nil {
		return errors.Wrap(c.trackAWSErrors(awserr), "getting cluster list")
	}

	collectorFn := c.metricStats
	if c.useGMD {
		collectorFn = c.metricData
	}
	var buf bytes.Buffer
	buf.Grow(32768)

	for cid, nodes := range clusterList {
		// GetMetricData and GetMetricStatistics both have their pros and cons...
		// let customer decide on a per-collector configuration basis which
		// is best for the use-case
		for _, nid := range nodes {
			dimensions := []*cloudwatch.Dimension{
				{
					Name:  aws.String("CacheClusterId"),
					Value: aws.String(cid),
				},
				{
					Name:  aws.String("CacheNodeId"),
					Value: aws.String(nid),
				},
			}
			if err := collectorFn(&buf, sess, timespan, dimensions, baseTags); err != nil {
				c.logger.Warn().Err(err).Str("cluster_id", cid).Str("node_id", nid).Msg("fetching telemetry")
				continue
			}
			if buf.Len() == 0 {
				c.logger.Warn().Str("collector", c.ID()).Msg("no telemetry to submit")
				continue
			}
			c.logger.Debug().Str("collector", c.ID()).Msg("submitting telemetry")
			if err := c.check.SubmitMetrics(&buf); err != nil {
				c.logger.Warn().Err(err).Str("cluster_id", cid).Str("node_id", nid).Msg("submitting telemetry")
			}
			buf.Reset()
		}
	}

	return nil
}

func (c *ElastiCache) clusterList(sess client.ConfigProvider) (map[string][]string, error) {
	clusterList := map[string][]string{}
	ecSvc := elasticache.New(sess)

	// Get details on cache clusters (will list all clusters if no specific
	// cluster(s) defined in user config). Ignore any not in the 'available'
	// state.

	c.logger.Debug().Msg("getting aws elasticache cluster info")
	if c.clusterIDs != nil {
		for _, id := range *c.clusterIDs {
			dcci := &elasticache.DescribeCacheClustersInput{
				CacheClusterId:    aws.String(id),
				ShowCacheNodeInfo: aws.Bool(true),
			}
			cl, err := ecSvc.DescribeCacheClusters(dcci)
			if err != nil {
				var awsErr awserr.Error
				if errors.As(err, &awsErr) {
					return nil, fmt.Errorf("describing cluster: %w", err)
				}
				c.logger.Warn().Err(err).Str("cluster_id", id).Msg("describing elasticache cluster, skipping")
				continue
			}
			for _, cluster := range cl.CacheClusters {
				if *cluster.CacheClusterStatus != "available" {
					c.logger.Debug().Str("cluter_id", *cluster.CacheClusterId).Str("status", *cluster.CacheClusterStatus).Msg("invalid state, skipping")
					continue
				}
				clusterList[*cluster.CacheClusterId] = []string{}
				for _, node := range cluster.CacheNodes {
					clusterList[*cluster.CacheClusterId] = append(clusterList[*cluster.CacheClusterId], *node.CacheNodeId)
				}
			}
		}

		return clusterList, nil
	}

	dcci := &elasticache.DescribeCacheClustersInput{
		ShowCacheNodeInfo: aws.Bool(true),
	}
	cl, err := ecSvc.DescribeCacheClusters(dcci)
	if err != nil {
		return nil, fmt.Errorf("describing elasticache clusters: %w", err)
	}
	for _, cluster := range cl.CacheClusters {
		if *cluster.CacheClusterStatus != "available" {
			c.logger.Debug().Str("cluter_id", *cluster.CacheClusterId).Str("status", *cluster.CacheClusterStatus).Msg("invalid state, skipping")
			continue
		}
		clusterList[*cluster.CacheClusterId] = []string{}
		for _, node := range cluster.CacheNodes {
			clusterList[*cluster.CacheClusterId] = append(clusterList[*cluster.CacheClusterId], *node.CacheNodeId)
		}
	}

	return clusterList, nil
}

// DefaultMetrics defines the default set of metrics for the service
// https://docs.aws.amazon.com/AmazonElastiCache/latest/mem-ug/CacheMetrics.HostLevel.html
// https://docs.aws.amazon.com/AmazonElastiCache/latest/mem-ug/CacheMetrics.Memcached.html
// https://docs.aws.amazon.com/AmazonElastiCache/latest/mem-ug//CacheMetrics.WhichShouldIMonitor.html
func (c *ElastiCache) DefaultMetrics() []Metric {
	return []Metric{
		{
			AWSMetric{
				Disabled: false,
				Name:     "CPUUtilization",
				Stats:    []string{metricStatAverage},
				Units:    "Percent",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "FreeableMemory",
				Stats:    []string{metricStatAverage},
				Units:    "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "NetworkBytesIn",
				Stats:    []string{metricStatAverage},
				Units:    "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "NetworkBytesOut",
				Stats:    []string{metricStatAverage},
				Units:    "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "NetworkPacketsIn",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically

			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "NetworkPacketsOut",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: false,
				Name:     "SwapUsage",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "BytesReadIntoMemcached",
				Stats:    []string{metricStatAverage},
				Units:    "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "BytesUsedForCacheItems",
				Stats:    []string{metricStatAverage},
				Units:    "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "BytesWrittenOutFromMemcached",
				Stats:    []string{metricStatAverage},
				Units:    "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "CasBadval",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "CasHits",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "CasMisses",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "CmdFlush",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "CmdGet",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "CmdSet",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: false,
				Name:     "CurrConnections",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "CurrItems",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "DecrHits",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "DecrMisses",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "DeleteHits",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "DeleteMisses",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: false,
				Name:     "Evictions",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "GetHits",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "GetMisses",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "IncrHits",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "IncrMisses",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "Reclaimed",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "BytesUsedForHash",
				Stats:    []string{metricStatAverage},
				Units:    "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "CmdConfigGet",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "CmdConfigSet",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "CmdTouch",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "CurrConfig",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "EvictedUnfetched",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "ExpiredUnfetched",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "SlabsMoved",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "TouchHits",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "TouchMisses",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "NewConnections",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "NewItems",
				Stats:    []string{metricStatAverage},
				Units:    "Count",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
		{
			AWSMetric{
				Disabled: true,
				Name:     "UnusedMemory",
				Stats:    []string{metricStatAverage},
				Units:    "Bytes",
			},
			CirconusMetric{
				Name: "",              // NOTE: AWSMetric.Name will be used if blank
				Type: "gauge",         // (gauge|counter|histogram|text)
				Tags: circonus.Tags{}, // NOTE: units:strings.ToLower(AWSMetric.Units) is added automatically
			},
		},
	}
}
