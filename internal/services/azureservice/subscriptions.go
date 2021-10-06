// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package azureservice

import (
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-06-01/subscriptions"
	"github.com/pkg/errors"
)

// subscriptionMeta is the meta data about a subscription used to create a Circonus check.
type subscriptionMeta struct {
	Name string
}

// getSubscriptionMeta retrieves the subscription meta data used to create a Circonus check.
func (inst *Instance) getSubscriptionMeta() (*subscriptionMeta, error) {
	auth, err := inst.authorize()
	if err != nil {
		return nil, errors.Wrap(err, "authorize, subscription meta")
	}

	subscriptionClient := subscriptions.NewClient()
	subscriptionClient.Authorizer = auth
	if err = subscriptionClient.AddToUserAgent(inst.cfg.Azure.UserAgent); err != nil {
		inst.logger.Warn().Err(err).Msg("adding user agent to client")
	}
	// ref: https://godoc.org/github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-06-01/subscriptions#Client.Get
	result, err := subscriptionClient.Get(inst.ctx, inst.cfg.Azure.SubscriptionID)
	if err != nil {
		return nil, err
	}

	if result.State != subscriptions.Enabled {
		return nil, errors.Errorf("subscription is not enabled - id: %s, name: %s, state: %s", *result.ID, *result.DisplayName, result.State)
	}

	sm := subscriptionMeta{
		Name: *result.DisplayName,
	}

	return &sm, nil
}
