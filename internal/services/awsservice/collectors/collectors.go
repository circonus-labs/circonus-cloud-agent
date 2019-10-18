// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	metricStatAverage     = "Average"
	metricStatSum         = "Sum"
	metricStatMaximum     = "Maximum"
	metricStatMinimum     = "Minimum"
	metricStatSampleCount = "SampleCount"
	resultIDFormat        = "m%ds%sq%d" //  config metric index, metric stat name, result metric counter index
)

// Collector interface for aws metric services
type Collector interface {
	// Collect(sess *session.Session, metricDest io.Writer, baseTags circonus.Tags, interval uint, period int64) error
	Collect(sess *session.Session, timespan MetricTimespan, baseTags circonus.Tags) error
	ID() string
	DefaultMetrics() []Metric
}

// AWSCollector defines a generic aws service metric collector
type AWSCollector struct {
	Namespace  string                  `json:"namespace" toml:"namespace" yaml:"namespace"`    // e.g. AWS/EC2
	Disabled   bool                    `json:"disabled" toml:"disabled" yaml:"disabled"`       // disable metric collection for this aws service namespace
	Dimensions []*cloudwatch.Dimension `json:"dimensions" toml:"dimensions" yaml:"dimensions"` // key:val pairs
	Metrics    []Metric                `json:"metrics" toml:"metrics" yaml:"metrics"`          // mapping of metrics to collect
	Tags       circonus.Tags           `json:"tags" toml:"tags" yaml:"tags"`                   // service tags
	UseGMD     bool                    `json:"use_gmd" toml:"use_gmd" yaml:"use_gmd"`          // use getMetricData instead of getMetricStatsistics
	// EC2 only
	InstanceFilters *[]Filter `json:"instance_filters,omitempty" toml:"instance_filters,omitempty" yaml:"instance_filters,omitempty"`
	// ElastiCache only
	CacheClusterIDs *[]string `json:"cache_cluster_ids,omitempty" toml:"cache_cluster_ids,omitempty" yaml:"cache_cluster_ids,omitempty"`
}

// AWSMetric defines an AWS metrics
type AWSMetric struct {
	Disabled bool     `json:"disabled" toml:"disabled" yaml:"disabled"` // disable collection (DEFAULT: false)
	Name     string   `json:"name" toml:"name" yaml:"name"`             // REQUIRED
	Stats    []string `json:"stats" toml:"stats" yaml:"stats"`          // REQUIRED
	Units    string   `json:"units" toml:"units" yaml:"units"`          // REQUIRED
}

// CirconusMetric defines a Circonus metric
type CirconusMetric struct {
	Name string        `json:"name" toml:"name" yaml:"name"` // DEFAULT AWSMetric.Name
	Type string        `json:"type" toml:"type" yaml:"type"` // REQUIRED, (gauge|counter|histogram|text) - NOTE: histogram is for HIGH volume data (as in multiple samples per second, aws does not provide this level of granularity)
	Tags circonus.Tags `json:"tags" toml:"tags" yaml:"tags"` // DEFAULT none - additional metric specific stream tags - will automatically add "units:"+strings.ToLower(AWSMetric.Units), aws_region:current_region_being_polled, aws_instance_id:instance_id_of_metric_origin
}

// Metric maps a given metric between AWS and Circonus
type Metric struct {
	AWSMetric      AWSMetric      `json:"aws" toml:"aws" yaml:"aws"`                // REQUIRED
	CirconusMetric CirconusMetric `json:"circonus" toml:"circonus" yaml:"circonus"` // REQUIRED
}

// Filter defines a generic AWS EC2 filter https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#Filter
type Filter struct {
	Name   *string   `json:"name" toml:"name" yaml:"name"`
	Values []*string `json:"values" toml:"values" yaml:"values"`
}

// MetricTimespan defines the span of time for requesting metrics
type MetricTimespan struct {
	Start  time.Time
	End    time.Time
	Period int64
}

// New creates a new collector instance
func New(ctx context.Context, check *circonus.Check, cfgs []AWSCollector, logger zerolog.Logger) ([]Collector, error) {
	// TODO: zone/service discovery (so that bare minimum config required would be credentials)
	//
	// aws cloudwatch call(s) for what services are in use (if no list all active services, call each for list metrics and use ones that return >0 metrics)
	// walk through active service list, activate ones we currently support unless explicitly disabled
	// create default service config
	//      do we have a custom service config
	//          if disabled, skip this service
	//          if not disabled, overlay custom config
	// add configured service to collectors slice

	cl := collectorList()
	cc := []Collector{}
	for _, cfg := range cfgs {
		cfg := cfg
		if cfg.Disabled {
			continue // allow metrics from an entire service to be disabled
		}
		var c Collector
		var err error

		if initfn, known := cl[strings.ToLower(cfg.Namespace)]; known {
			c, err = initfn(ctx, check, &cfg, logger)
		} else {
			err = errors.New("unrecognized aws service namespace")
		}

		if err != nil {
			logger.Warn().Err(err).Str("namespace", cfg.Namespace).Msg("skipping")
			continue
		}

		if c == nil {
			logger.Warn().Err(errors.New("init returned nil collector")).Str("namespace", cfg.Namespace).Msg("skipping")
			continue
		}

		cc = append(cc, c)
	}

	if len(cc) == 0 {
		return nil, errors.New("no services configured from which to collect metrics")
	}

	return cc, nil
}

type collectorInitFn func(context.Context, *circonus.Check, *AWSCollector, zerolog.Logger) (Collector, error)
type collectorInitList map[string]collectorInitFn

func collectorList() collectorInitList {
	return collectorInitList{
		"aws/applicationelb":    newApplicationELB,
		"aws/cloudfront":        newCloudFront,
		"aws/dynamodb":          newDynamoDB,
		"aws/dx":                newDX,
		"aws/ebs":               newEBS,
		"aws/ec2":               newEC2,
		"aws/ec2autoscaling":    newEC2AutoScaling,
		"aws/ec2spot":           newEC2Spot,
		"aws/ecs":               newECS,
		"aws/efs":               newEFS,
		"aws/elasticbeanstalk":  newElasticBeanstalk,
		"aws/elasticache":       newElastiCache,
		"aws/elasticinterface":  newElasticInterface,
		"aws/elasticmapreduce":  newElasticMapReduce,
		"aws/elastictranscoder": newElasticTranscoder,
		"aws/elb":               newELB,
		"aws/es":                newES,
		"aws/kms":               newKMS,
		"aws/lambda":            newLambda,
		"aws/networkelb":        newNetworkELB,
		"aws/rds":               newRDS,
		"aws/route53":           newRoute53,
		"aws/route53resolver":   newRoute53Resolver,
		"aws/s3":                newS3,
		"aws/sns":               newSNS,
		"aws/sqs":               newSQS,
		"aws/natgateway":        newNATGateway,     // VPC
		"aws/transitgateway":    newTransitGateway, // VPC
	}
}

// ConfigExample generates configuration examples for collectors
// nolint: gocyclo
func ConfigExample() ([]AWSCollector, error) {
	// NOTE: Certain services (e.g. aws/applicationelb *require* dimensions,
	//       there are no default metrics w/o dimensions.) A blank configuration
	//       will be emitted because it requires additional user-supplied info.
	var cc []AWSCollector
	cl := collectorList()
	for cn := range cl {
		c := AWSCollector{
			Namespace: cn,
			Disabled:  true,
		}
		var v Collector
		switch cn {
		case "aws/applicationelb":
			v = &ApplicationELB{}
		case "aws/cloudfront":
			v = &CloudFront{}
		case "aws/dynamodb":
			v = &DynamoDB{}
		case "aws/dx":
			v = &DX{}
		case "aws/ebs":
			v = &EBS{}
		case "aws/ec2":
			c.Disabled = false
			c.InstanceFilters = &[]Filter{}
			v = &EC2{}
		case "aws/ec2autoscaling":
			v = &EC2AutoScaling{}
		case "aws/ec2spot":
			v = &EC2Spot{}
		case "aws/ecs":
			v = &ECS{}
		case "aws/efs":
			v = &EFS{}
		case "aws/elasticbeanstalk":
			v = &ElasticBeanstalk{}
		case "aws/elasticache":
			v = &ElastiCache{}
		case "aws/elasticinterface":
			v = &ElasticInterface{}
		case "aws/elasticmapreduce":
			v = &ElasticMapReduce{}
		case "aws/elastictranscoder":
			v = &ElasticTranscoder{}
		case "aws/elb":
			v = &ELB{}
		case "aws/es":
			v = &ES{}
		case "aws/kms":
			v = &KMS{}
		case "aws/lambda":
			v = &Lambda{}
		case "aws/networkelb":
			v = &NetworkELB{}
		case "aws/rds":
			v = &RDS{}
		case "aws/route53":
			v = &Route53{}
		case "aws/route53resolver":
			v = &Route53Resolver{}
		case "aws/s3":
			v = &S3{}
		case "aws/sns":
			v = &SNS{}
		case "aws/sqs":
			v = &SQS{}
		case "aws/natgateway": // VPC
			v = &NATGateway{}
		case "aws/transitgateway": // VPC
			v = &TransitGateway{}
		default:
			return nil, errors.Errorf("unknown aws service namespace (%s)", cn)
		}
		if v != nil {
			c.Metrics = v.DefaultMetrics()
			cc = append(cc, c)
		}
	}

	return cc, nil
}

// nolint: structcheck
type common struct {
	id           string
	check        *circonus.Check
	enabled      bool
	disableCause string    // cause of a runtime disabling of the collector
	disableTime  time.Time // time of runtime disabling (will try again every hour)
	ctx          context.Context
	dimensions   []*cloudwatch.Dimension
	metrics      []Metric
	tags         circonus.Tags
	useGMD       bool
	logger       zerolog.Logger
}

func newCommon(ctx context.Context, ns string, check *circonus.Check, cfg *AWSCollector, logger zerolog.Logger) common {
	return common{
		id:         ns,
		enabled:    true,
		ctx:        ctx,
		check:      check,
		dimensions: cfg.Dimensions,
		metrics:    cfg.Metrics,
		tags:       cfg.Tags,
		useGMD:     cfg.UseGMD,
		logger:     logger.With().Str("collector", ns).Logger(),
	}
}

// Collect is the common collection method - can be overridden by individual collectors (see ec2/elasticache)
func (c *common) Collect(sess *session.Session, timespan MetricTimespan, baseTags circonus.Tags) error {
	if sess == nil {
		return errors.New("invalid session (nil)")
	}

	if !c.Enabled() {
		return nil
	}

	// GetMetricData and GetMetricStatistics both have their pros and cons...
	// let customer decide on a per-config basis which is best for the use-case
	collectorFn := c.metricStats
	if c.useGMD {
		collectorFn = c.metricData
	}

	var buf bytes.Buffer
	buf.Grow(32768)

	c.logger.Debug().Str("collector", c.ID()).Msg("collecting telemetry")
	var dims []*cloudwatch.Dimension
	if err := collectorFn(&buf, sess, timespan, dims, baseTags); err != nil {
		return errors.Wrap(err, "collecting telemetry")
	}

	if buf.Len() == 0 {
		c.logger.Warn().Str("collector", c.ID()).Msg("no telemetry to submit")
		return nil
	}

	c.logger.Debug().Str("collector", c.ID()).Msg("submitting telemetry")
	if err := c.check.SubmitMetrics(&buf); err != nil {
		return errors.Wrap(err, "submitting telemetry")
	}

	return nil
}

// Enabled returns true if the collector is enabled. If the collector has been
// dynamically disabled by a potentially temporary error, that will be logged
// and the collector will be re-enabled to try again after an hour. See the
// AccessDenied case in trackAWSErrors method.
func (c *common) Enabled() bool {
	if c.enabled {
		return c.enabled
	}

	// configured to be disabled, do not emit log message
	if c.disableCause == "" {
		return false
	}

	// disabled by a _potentially_ intermittent error (e.g. auth snafu)
	if time.Since(c.disableTime) >= 1*time.Hour { // TODO: TBD allow disable retry time to be configurable
		c.enabled = true // re-enable and try again
		return c.enabled
	}

	c.logger.Warn().Time("disable_time", c.disableTime).Str("disable_cause", c.disableCause).Msg("collector has been disabled")
	return false
}

// ID returns the colletor's string id/name (e.g. used in logging by the instance using the collector)
func (c *common) ID() string {
	return c.id
}

// done returns true if the context has been cancelled
func (c *common) done() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
		return false
	}
}

// trackAWSErrors will unravel errors and log the specific aws errors. it will return
// the specific error if it is an aws error otherwise it returns the original error
// or nil if there was no error passed.
func (c *common) trackAWSErrors(err error) error {
	if err == nil {
		return err
	}

	if awsErr, ok := errors.Cause(err).(awserr.Error); ok {
		if reqErr, ok := errors.Cause(err).(awserr.RequestFailure); ok {
			c.logger.Error().
				Str("aws_req_code", reqErr.Code()).
				Str("aws_req_msg", reqErr.Message()).
				Str("req_id", reqErr.RequestID()).
				Msg("aws request error")
			if reqErr.Code() == "AccessDenied" { // AccessDenied to resource
				c.enabled = false
				c.disableTime = time.Now()
				c.disableCause = reqErr.Message()
				return reqErr
			}
		} else {
			c.logger.Error().
				Str("aws_code", awsErr.Code()).
				Str("aws_msg", awsErr.Message()).
				Msg("aws error")
			return awsErr
		}
	}

	return err
}
