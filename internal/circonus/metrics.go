// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package circonus

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/circonus-labs/circonus-cloud-agent/internal/release"
	"github.com/pkg/errors"
)

// SubmitMetrics to Circonus check.
func (c *Check) SubmitMetrics(metricSrc io.Reader) error {
	c.Lock()
	defer c.Unlock()

	if c.bundle == nil {
		return errors.New("invalid state (nil check bundle)")
	}

	if metricSrc == nil {
		return errors.New("invalid metric source (nil)")
	}

	subURL := ""
	if surl, found := c.bundle.Config["submission_url"]; found {
		subURL = surl
	} else {
		return errors.New("invalid check bundle, no submission url")
	}

	var client *http.Client

	if c.brokerTLS != nil {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				TLSClientConfig:     c.brokerTLS,
				TLSHandshakeTimeout: 10 * time.Second,
				DisableKeepAlives:   false,
				MaxIdleConnsPerHost: 2,
				DisableCompression:  false,
			},
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				DisableKeepAlives:   false,
				MaxIdleConnsPerHost: 2,
				DisableCompression:  false,
			},
		}
	}

	// for debugging only
	mbuff, err := ioutil.ReadAll(metricSrc)
	if err != nil {
		return err
	}
	if e := c.logger.Debug(); e.Enabled() {
		if c.config.TraceMetrics {
			e.Msg("Submitted data")
			fmt.Printf("\n===BEGIN(%d)\n%s\n===END\n", time.Now().UTC().UnixNano(), mbuff)
		}
	}
	req, err := http.NewRequestWithContext(context.Background(), "PUT", subURL, bytes.NewReader(mbuff))
	// return to this one when debugging submissions is complete
	// req, err := http.NewRequest("PUT", subURL, metricSrc)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", release.NAME+"/"+release.VERSION)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Connection", "close")

	resp, err := client.Do(req)
	if err != nil {
		client.CloseIdleConnections()
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close() // nolint: errcheck
	if err != nil {
		client.CloseIdleConnections()
		return err
	}

	if resp.StatusCode != http.StatusOK {
		if pr, isPipeReader := metricSrc.(*io.PipeReader); isPipeReader {
			if err = pr.Close(); err != nil {
				c.logger.Warn().Err(err).Msg("closing pipe reader")
			}
		}
		c.logger.Error().Err(err).Str("url", subURL).Str("status", resp.Status).RawJSON("response", body).Msg("submitting telemetry")
		client.CloseIdleConnections()
		return errors.Wrap(err, "submitting metrics")
	}

	c.logger.Debug().Str("cid", c.bundle.CID).RawJSON("result", body).Msg("telmetry stats submitted")

	client.CloseIdleConnections()

	return nil
}

// WriteMetricSample to queue for submission.
func (c *Check) WriteMetricSample(metricDest io.Writer, metricName, metricType string, value interface{}, timestamp *time.Time) error {
	if metricDest == nil {
		return errors.New("invalid metric destination (nil)")
	}
	if metricName == "" {
		return errors.New("invalid metric name (empty)")
	}
	if metricType == "" {
		return errors.New("invalid metric type (empty)")
	}

	if len(metricName) > MaxMetricNameLen {
		c.logger.Warn().
			Str("metric_name", metricName).
			Int("encoded_len", len(metricName)).
			Int("max_len", MaxMetricNameLen).
			Msg("max metric name length exceeded, discarding")
		return nil
	}

	if !c.metricTypeRx.MatchString(metricType) {
		return errors.Errorf("unrecognized circonus metric type (%s)", metricType)
	}

	val := value
	if metricType == "s" {
		val = fmt.Sprintf(`"%v"`, strings.ReplaceAll(value.(string), `"`, `\"`)) // NOTE: need to insert the quotes
	}

	// escape quotes in streamtags
	metricName = strings.ReplaceAll(metricName, `"`, `\"`)

	var metricSample string
	if timestamp != nil {
		metricSample = fmt.Sprintf(
			`{"%s":{"_type":"%s","_value":%v,"_ts":%d}}`,
			metricName,
			metricType,
			val,
			uint64(timestamp.UTC().Unix()*1000), // trap wants milliseconds
		)
	} else {
		metricSample = fmt.Sprintf(
			`{"%s":{"_type":"%s","_value":%v}}`,
			metricName,
			metricType,
			val,
		)
	}

	if c.config.TraceMetrics {
		c.logger.Debug().Str("metric", metricSample).Msg("writing")
	}

	_, err := fmt.Fprintln(metricDest, metricSample)
	if err != nil {
		return err
	}
	return nil
}
