// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package circonus

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	apiclient "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
)

// initializeBroker fetches broker from circonus api and sets up broker tls config.
func (c *Check) initializeBroker() error {
	if c.apih == nil {
		return errors.New("invalid state (nil api client)")
	}
	if c.bundle == nil {
		return errors.New("invalid state (nil bundle)")
	}

	if len(c.bundle.Brokers) == 0 {
		return errors.New("invalid bundle, 0 brokers")
	}

	cid := c.bundle.Brokers[0]
	broker, err := c.apih.FetchBroker(apiclient.CIDType(&cid))
	if err != nil {
		return errors.Wrap(err, "fetching broker")
	}

	c.broker = broker

	return nil
}

// BrokerTLSConfig returns the broker tls configuration for metric submissions or nil if tls config not needed (e.g. public trap broker).
func (c *Check) BrokerTLSConfig() (*tls.Config, error) {
	c.Lock()
	defer c.Unlock()

	if c.broker == nil {
		return nil, errors.New("invalid state (nil broker)")
	}

	if strings.Contains(c.bundle.Config["submission_url"], "api.circonus.com") {
		return nil, nil // api.circonus.com uses a public certificate, no tls config needed
	}

	if c.brokerTLS != nil {
		return c.brokerTLS, nil
	}

	if err := c.setBrokerTLSConfig(); err != nil {
		return nil, errors.Wrap(err, "setting broker tls config")
	}

	return c.brokerTLS, nil
}

// brokerCN returns broker cn based on broker object.
func (c *Check) brokerCN(submissionURL string) (string, error) {
	if c.broker == nil {
		return "", errors.New("invalid state (nil broker)")
	}
	u, err := url.Parse(submissionURL)
	if err != nil {
		return "", errors.Wrap(err, "determining broker cn")
	}

	hostParts := strings.Split(u.Host, ":")
	host := hostParts[0]

	if net.ParseIP(host) == nil { // it's a non-ip string
		return u.Host, nil
	}

	cn := ""

	for _, detail := range c.broker.Details {
		if *detail.IP == host {
			cn = detail.CN
			break
		}
	}

	if cn == "" {
		return "", errors.Errorf("error, unable to match URL host (%s) to Broker", u.Host)
	}

	return cn, nil
}

// setBrokerTLSConfig sets up the broker tls config for metric submissions.
func (c *Check) setBrokerTLSConfig() error {
	cn, err := c.brokerCN(c.bundle.Config["submission_url"])
	if err != nil {
		return fmt.Errorf("unable to determine broker CN: %w", err)
	}

	if c.config.BrokerCAFile != "" {
		cert, err := os.ReadFile(c.config.BrokerCAFile) //nolint:govet
		if err != nil {
			return errors.Wrap(err, "configuring broker tls")
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(cert) {
			return errors.New("unable to add Broker CA Certificate to x509 cert pool")
		}
		c.brokerTLS = &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    cp,
			ServerName: cn,
		}
		return nil
	}

	type cacert struct {
		Contents string `json:"contents"`
	}

	jsoncert, err := c.apih.Get("/pki/ca.crt")
	if err != nil {
		return errors.Wrap(err, "fetching broker ca cert from api")
	}
	var cadata cacert
	if err := json.Unmarshal(jsoncert, &cadata); err != nil {
		return errors.Wrap(err, "parsing broker ca cert from api")
	}
	if cadata.Contents == "" {
		return errors.Errorf("unable to find ca cert 'Contents' attribute in api response (%+v)", cadata)
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM([]byte(cadata.Contents)) {
		return errors.New("unable to add Broker CA Certificate to x509 cert pool")
	}
	c.brokerTLS = &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    cp,
		ServerName: cn,
	}
	return nil
}
