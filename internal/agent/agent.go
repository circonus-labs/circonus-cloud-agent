// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agent

import (
	"context"
	"os"
	"os/signal"

	"github.com/circonus-labs/circonus-cloud-agent/internal/config"
	"github.com/circonus-labs/circonus-cloud-agent/internal/release"
	"github.com/circonus-labs/circonus-cloud-agent/internal/services"
	"github.com/circonus-labs/circonus-cloud-agent/internal/services/awsservice"
	"github.com/circonus-labs/circonus-cloud-agent/internal/services/azureservice"
	"github.com/circonus-labs/circonus-cloud-agent/internal/services/gcpservice"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// Agent holds the main circonus-agent process.
type Agent struct {
	group       *errgroup.Group
	groupCtx    context.Context
	groupCancel context.CancelFunc
	services    map[string]services.Service
	signalCh    chan os.Signal
}

// New returns a new agent instance.
func New() (*Agent, error) {
	ctx, cancel := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	var err error
	a := Agent{
		group:       g,
		groupCtx:    gctx,
		groupCancel: cancel,
		services:    make(map[string]services.Service),
		signalCh:    make(chan os.Signal, 10),
	}

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	{ // AWS
		awssvc, err := awsservice.New(a.groupCtx)
		if err != nil {
			return nil, errors.Wrap(err, "creating AWS client")
		}
		if awssvc.Enabled() {
			a.services["aws"] = awssvc
		}
	}

	{ // Azure
		azuresvc, err := azureservice.New(a.groupCtx)
		if err != nil {
			return nil, errors.Wrap(err, "creating Azure client")
		}
		if azuresvc.Enabled() {
			a.services["azure"] = azuresvc
		}
	}

	{ // GCP
		gcpsvc, err := gcpservice.New(a.groupCtx)
		if err != nil {
			return nil, errors.Wrap(err, "creating GCP client")
		}
		if gcpsvc.Enabled() {
			a.services["gcp"] = gcpsvc
		}
	}

	if len(a.services) == 0 {
		log.Fatal().Msg("no cloud services enabled, must enable at least ONE")
	}

	a.signalNotifySetup()

	return &a, nil
}

// Start the agent.
func (a *Agent) Start() error {
	a.group.Go(a.handleSignals)
	for svcID := range a.services {
		a.group.Go(a.services[svcID].Start)
	}

	log.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("starting wait")

	return a.group.Wait()
}

// Stop cleans up and shuts down the Agent.
func (a *Agent) Stop() {
	a.stopSignalHandler()
	a.groupCancel()

	log.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("stopped")
}

// stopSignalHandler disables the signal handler.
func (a *Agent) stopSignalHandler() {
	signal.Stop(a.signalCh)
	signal.Reset() // so a second ctrl-c will force immediate stop
}
