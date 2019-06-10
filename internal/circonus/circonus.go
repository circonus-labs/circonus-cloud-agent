// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package circonus

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"

	"github.com/circonus-labs/circonus-cloud-agent/internal/release"
	apiclient "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// ServiceConfig defines the Circonus configuration for a cloud service to include
type ServiceConfig struct {
	CID          string `json:"cid" toml:"cid" yaml:"cid"`                               // OPTIONAL, cid of specific check bundle to use
	BrokerCID    string `json:"broker_cid" toml:"broker_cid" yaml:"broker_cid"`          // OPTIONAL, default public trap broker
	App          string `json:"app" toml:"app" yaml:"app"`                               // REQUIRED
	Key          string `json:"key" toml:"key" yaml:"key"`                               // REQUIRED
	URL          string `json:"url" toml:"url" yaml:"url"`                               // DEFAULT 'https://api.circonus.com/v2/'
	CAFile       string `json:"ca_file" toml:"ca_file" yaml:"ca_file"`                   // DEFAULT api.circonus.com uses a public certificate
	Debug        bool   `json:"debug" toml:"debug" yaml:"debug"`                         // DEFAULT false - this is separate so that the global debug does not innundate logs with cgm debug messages from each cloud service metric collection client
	TraceMetrics bool   `json:"trace_metrics" toml:"trace_metrics" yaml:"trace_metrics"` // DEFAULT false - output each metric as it is sent
}

// Config options for a Circonus Check passed to New method
type Config struct {
	ID            string         // a unique identifier, used to search for a check bundle
	CheckBundleID string         // a specific check bundle cid to use
	BrokerCID     string         // broker cid to use (default circonus public broker)
	BrokerCAFile  string         // broker ca file
	DisplayName   string         // display name to use for check bundle (when searching or creating)
	Tags          string         // tags to add to a check when creating
	APIKey        string         // Circonus API Token key
	APIApp        string         // Circonus API Token app
	APIURL        string         // Circonus API URL
	APICAFile     string         // api ca file
	Debug         bool           // turn on debugging messages
	Logger        zerolog.Logger // logging instance to use
	TraceMetrics  bool           // output each metric as it is sent
}

// Check defines a Circonus check for a circonus-cloud-agent service
type Check struct {
	sync.Mutex
	apih            *apiclient.API
	config          *Config
	errorMetricName string
	broker          *apiclient.Broker
	brokerTLS       *tls.Config
	bundle          *apiclient.CheckBundle
	metricTypeRx    *regexp.Regexp // validate metric types
	logger          zerolog.Logger
}

// logshim is used to satisfy apiclient Logger interface (avoiding ptr receiver issue)
type logshim struct {
	logh zerolog.Logger
}

func (l logshim) Printf(fmt string, v ...interface{}) {
	l.logh.Printf(fmt, v...)
}

const (
	// MetricTypeInt32 reconnoiter
	MetricTypeInt32 = "i"

	// MetricTypeUint32 reconnoiter
	MetricTypeUint32 = "I"

	// MetricTypeInt64 reconnoiter
	MetricTypeInt64 = "l"

	// MetricTypeUint64 reconnoiter
	MetricTypeUint64 = "L"

	// MetricTypeFloat64 reconnoiter
	MetricTypeFloat64 = "n"

	// MetricTypeString reconnoiter
	MetricTypeString = "s"

	// NOTE: max tags and metric name len are enforced here so that
	// details on which metric(s) can be logged. Otherwise, any
	// metric(s) exceeding the limits are rejected by the broker
	// without details on exactly which metric(s) caused the error.
	// All metrics sent with the offending metric(s) are also rejected.

	// MaxTags reconnoiter will accept in stream tagged metric name
	MaxTags = 256 // sync w/MAX_TAGS https://github.com/circonus-labs/reconnoiter/blob/master/src/noit_metric.h#L41

	// MaxMetricNameLen reconnoiter will accept (name+stream tags)
	MaxMetricNameLen = 4096 // sync w/MAX_METRIC_TAGGED_NAME https://github.com/circonus-labs/reconnoiter/blob/master/src/noit_metric.h#L40
)

var (
	publicHTTPTrapBrokerCID = "/broker/35"
	checkType               = "httptrap"
	checkStatusActive       = "active"
	checkMetricFilters      = [][]string{
		{"deny", "^$", ""},
		{"allow", "^.+$", ""},
	}
)

// NewCheck creates a new Circonus check instance based on the Config options passed to
// initialize the Circonus API, check and broker.
func NewCheck(cfg *Config) (*Check, error) {
	if cfg == nil {
		return nil, errors.New("invalid config (nil)")
	}

	c := &Check{
		config:          cfg,
		errorMetricName: strings.ReplaceAll(release.NAME, "-", "_") + "_errors", // TODO: may become a config option
		logger:          cfg.Logger.With().Str("pkg", "check").Logger(),
		metricTypeRx: regexp.MustCompile("^[" + strings.Join([]string{
			MetricTypeInt32,
			MetricTypeUint32,
			MetricTypeInt64,
			MetricTypeUint64,
			MetricTypeFloat64,
			MetricTypeString,
		}, "") + "]$"),
	}

	if err := c.initAPI(); err != nil {
		return nil, errors.Wrap(err, "initializing Circonus API")
	}

	if err := c.initializeCheckBundle(); err != nil {
		return nil, errors.Wrap(err, "initializing check")
	}

	if err := c.initializeBroker(); err != nil {
		return nil, errors.Wrap(err, "initializing broker")
	}

	return c, nil
}

// RefreshCheck fetches a new copy of the check bundle and broker from Circonus API.
// Primary use-case is when a check is moved from one broker to another, metric
// submits will fail. Refreshing the check bundle will obtain the new submission
// url, and reconfigure the broker tls config (if needed).
func (c *Check) RefreshCheck() error {
	c.Lock()
	defer c.Unlock()

	if c.apih == nil {
		return errors.New("invalid state (nil api client)")
	}

	c.bundle = nil
	c.broker = nil
	c.brokerTLS = nil

	if err := c.initializeCheckBundle(); err != nil {
		return errors.Wrap(err, "refreshing check")
	}

	if err := c.initializeBroker(); err != nil {
		return errors.Wrap(err, "refreshing broker")
	}

	return nil
}

// initAPI creates and configures a Circonus API client
func (c *Check) initAPI() error {
	c.logger.Debug().Msg("initializing api client")
	apiConfig := &apiclient.Config{
		TokenKey: c.config.APIKey,
		TokenApp: c.config.APIApp,
		URL:      c.config.APIURL,
		Debug:    c.config.Debug,
		Log:      logshim{logh: c.logger.With().Str("pkg", "apicli").Logger()},
	}
	if c.config.APICAFile != "" {
		cert, err := ioutil.ReadFile(c.config.APICAFile)
		if err != nil {
			return errors.Wrap(err, "configuring API client")
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(cert) {
			return errors.New("unable to add API CA Certificate to x509 cert pool")
		}
		apiConfig.TLSConfig = &tls.Config{RootCAs: cp}
	}
	client, err := apiclient.New(apiConfig)
	if err != nil {
		return errors.Wrap(err, "creating API client")
	}

	c.apih = client
	return nil
}
