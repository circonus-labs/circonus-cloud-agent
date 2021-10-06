// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package azureservice

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/pkg/errors"
)

// resourceMeta contains details on a resource for handling the metrics available.
type resourceMeta struct {
	ID   string
	Name string
	Kind string
	Type string
	Tags circonus.Tags
}

// getResources retrieves a list of resources in a subscription based on the user supplied filter.
func (inst *Instance) getResources(auth autorest.Authorizer) ([]resourceMeta, error) {
	resourceClient := resources.NewClient(inst.cfg.Azure.SubscriptionID)
	resourceClient.Authorizer = auth
	if err := resourceClient.AddToUserAgent(inst.cfg.Azure.UserAgent); err != nil {
		inst.logger.Warn().Err(err).Msg("adding user agent to client")
	}
	// ref: https://godoc.org/github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources#Client.ListComplete
	result, err := resourceClient.ListComplete(inst.ctx, inst.cfg.Azure.ResourceFilter, "", nil)
	if err != nil {
		return nil, err
	}

	rl := []resourceMeta{}
	for result.NotDone() {

		if err := result.NextWithContext(inst.ctx); err != nil {
			return nil, err
		}

		// ref: https://godoc.org/github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources#GenericResource
		v := result.Value()
		if v.ID != nil {
			// elide any resources that do not have metrics
			if ok, err := inst.resourceHasMetrics(auth, *v.ID); !ok {
				inst.logger.Debug().Str("resource_id", *v.ID).Msg("does not support metrics, skipping")
				continue
			} else if err != nil {
				inst.logger.Warn().Err(err).Str("resource_id", *v.ID).Msg("fetching metric namespaces, skipping")
				continue
			}

			res := resourceMeta{ID: *v.ID}

			var tags circonus.Tags
			if v.Name != nil {
				res.Name = *v.Name
				tags = append(tags, circonus.Tag{Category: "resource_name", Value: *v.Name})
			}
			if v.Kind != nil {
				res.Kind = *v.Kind
				tags = append(tags, circonus.Tag{Category: "resource_kind", Value: *v.Kind})
			}
			if v.Type != nil {
				res.Type = *v.Type
				tags = append(tags, circonus.Tag{Category: "resource_type", Value: *v.Type})
			}
			if v.Location != nil {
				tags = append(tags, circonus.Tag{Category: "resource_location", Value: *v.Location})
			}
			for cat, val := range v.Tags {
				tags = append(tags, circonus.Tag{Category: cat, Value: *val})
			}
			res.Tags = tags

			rl = append(rl, res)
		}

		if inst.done() {
			return rl, nil
		}
	}

	if len(rl) == 0 {
		return nil, errors.New("zero resources found")
	}

	return rl, nil
}

// resourceHasMetrics determines if a given resource has any metric namespaces.
// Not all resources expose metrics. A request to list metrics will return an
// error if the resource does not have any metrics. This call simply returns
// an empty list of metric namespaces which will make error handling of
// list metrics more straight-forward.
func (inst *Instance) resourceHasMetrics(auth autorest.Authorizer, resourceID string) (bool, error) {
	metricsNamespaceClient := insights.NewMetricNamespacesClient(inst.cfg.Azure.SubscriptionID)
	metricsNamespaceClient.Authorizer = auth
	if err := metricsNamespaceClient.AddToUserAgent(inst.cfg.Azure.UserAgent); err != nil {
		inst.logger.Warn().Err(err).Msg("adding user agent to client")
	}
	ts := time.Now().Add(-5 * time.Hour)
	// ref: https://godoc.org/github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights#MetricNamespacesClient.List
	result, err := metricsNamespaceClient.List(inst.ctx, resourceID, ts.Format(time.RFC3339))
	if err != nil {
		return false, err
	}

	return len(*result.Value) > 0, nil
}
