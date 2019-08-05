// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package circonus

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/circonus-labs/circonus-cloud-agent/internal/release"
	apiclient "github.com/circonus-labs/go-apiclient"
	apiclicfg "github.com/circonus-labs/go-apiclient/config"
	"github.com/pkg/errors"
)

// initializeCheckBundle finds or creates a new check bundle
func (c *Check) initializeCheckBundle() error {
	if c.apih == nil {
		return errors.New("invalid state (nil api client)")
	}

	cid := c.config.CheckBundleID

	if cid != "" {
		bundle, err := c.apih.FetchCheckBundle(apiclient.CIDType(&cid))
		if err != nil {
			return errors.Wrap(err, "fetching configured check bundle")
		}
		if bundle.Status != "active" {
			return errors.Errorf("invalid check bundle (%s), not active", bundle.CID)
		}

		c.bundle = bundle
		return nil
	}

	bundle, err := c.findOrCreateCheckBundle()
	if err != nil {
		return errors.Wrap(err, "finding/creating check")
	}

	c.logger.Debug().Interface("check_bundle", bundle).Msg("using check bundle")
	c.bundle = bundle

	return nil
}

// findOrCreateCheckBundle searches for a check bundle based on target and display name
func (c *Check) findOrCreateCheckBundle() (*apiclient.CheckBundle, error) {
	searchCriteria := apiclient.SearchQueryType(fmt.Sprintf(`(active:1)(type:"%s")(host:%s)`, checkType, c.config.ID))

	bundles, err := c.apih.SearchCheckBundles(&searchCriteria, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "searching for check (%s)", searchCriteria)
	}

	if len(*bundles) == 0 {
		return c.createCheckBundle()
	}

	numActive := 0
	checkIdx := -1
	for idx, cb := range *bundles {
		if cb.Status != checkStatusActive {
			continue
		}
		numActive++
		if checkIdx == -1 {
			checkIdx = idx // first match
		}
	}

	if numActive > 1 {
		return nil, errors.Errorf("multiple active checks found (%d) matching (%s)", numActive, searchCriteria)
	}

	bundle := (*bundles)[checkIdx]
	return &bundle, nil
}

// createCheckBundle creates a new check bundle
func (c *Check) createCheckBundle() (*apiclient.CheckBundle, error) {
	secret, err := makeSecret()
	if err != nil {
		secret = "myS3cr3t"
	}
	notes := fmt.Sprintf("%s-%s", release.NAME, release.VERSION)
	broker := c.config.BrokerCID
	if broker == "" {
		broker = publicHTTPTrapBrokerCID
	}

	checkConfig := &apiclient.CheckBundle{
		Brokers: []string{broker},
		Config: apiclient.CheckBundleConfig{
			"asynch_metrics": "true",
			"secret":         secret,
		},
		DisplayName:   c.config.DisplayName,
		MetricFilters: checkMetricFilters,
		MetricLimit:   apiclicfg.DefaultCheckBundleMetricLimit,
		Metrics:       []apiclient.CheckBundleMetric{},
		Notes:         &notes,
		Period:        60,
		Status:        checkStatusActive,
		Tags:          strings.Split(c.config.Tags, ","),
		Target:        c.config.ID,
		Timeout:       10,
		Type:          checkType,
	}

	bundle, err := c.apih.CreateCheckBundle(checkConfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating check")
	}

	return bundle, nil
}

func makeSecret() (string, error) {
	hash := sha256.New()
	x := make([]byte, 2048)
	if _, err := rand.Read(x); err != nil {
		return "", err
	}
	if _, err := hash.Write(x); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil))[0:16], nil
}
