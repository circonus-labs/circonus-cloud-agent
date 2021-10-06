// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package azureservice

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/pkg/errors"
)

// authorize creates the authorization token to use when making api calls.
func (inst *Instance) authorize() (autorest.Authorizer, error) {
	azureEnv, err := azure.EnvironmentFromName(inst.cfg.Azure.CloudName)
	if err != nil {
		return nil, errors.Wrap(err, "getting azure environment")
	}

	oauthConfig, err := adal.NewOAuthConfig(azureEnv.ActiveDirectoryEndpoint, inst.cfg.Azure.DirectoryID)
	if err != nil {
		return nil, errors.Wrap(err, "new oauth config")
	}

	token, err := adal.NewServicePrincipalToken(*oauthConfig, inst.cfg.Azure.ApplicationID, inst.cfg.Azure.ApplicationSecret, azureEnv.ResourceManagerEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "new service principal token")
	}

	return autorest.NewBearerAuthorizer(token), nil
}
