// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package gcpservice

import (
	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

// loadProjectMeta retrieves project meta data from gcp api
// https://godoc.org/google.golang.org/api/cloudresourcemanager/v1#Project
func (inst *Instance) loadProjectMeta() error {
	conf, err := google.JWTConfigFromJSON(inst.cfg.GCP.credentialData, "https://www.googleapis.com/auth/cloud-platform.read-only")
	if err != nil {
		return errors.Wrap(err, "initializing gcp credentials")
	}
	client := conf.Client(inst.ctx)

	// svc, err := cloudresourcemanager.New(client)
	svc, err := cloudresourcemanager.NewService(inst.ctx, option.WithHTTPClient(client))
	if err != nil {
		return errors.Wrap(err, "initializing resource manager service")
	}

	resp, err := svc.Projects.Get(inst.cfg.GCP.projectID).Context(inst.ctx).Do()
	if err != nil {
		return errors.Wrap(err, "retrieving project meta data")
	}

	// ensure project is active
	if resp.LifecycleState != "ACTIVE" {
		return errors.Errorf("project %s is not active (%s)", resp.Name, resp.LifecycleState)
	}

	// save the project name (for use in check bundle display name)
	inst.cfg.GCP.projectName = resp.Name

	// use project labels as base tags
	inst.baseTags = circonus.Tags{
		circonus.Tag{Category: "project_id", Value: inst.cfg.GCP.projectID},
	}
	for k, v := range resp.Labels {
		inst.baseTags = append(inst.baseTags, circonus.Tag{Category: k, Value: v})
	}

	return nil
}
