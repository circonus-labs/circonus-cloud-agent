// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"context"
	"strings"
	"time"

	"github.com/circonus-labs/circonus-cloud-agent/internal/circonus"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Collector interface for gcp metric services.
type Collector interface {
	Collect(timeseriesStart, timeseriesEnd time.Time, projectID string, creds []byte, baseTags circonus.Tags) error
	ID() string
}

// GCPCollector defines a generic gcp service metric collector.
type GCPCollector struct {
	Filter   Filter        `json:"filter" toml:"filter" yaml:"filter"`       // filter
	Name     string        `json:"name" toml:"name" yaml:"name"`             // e.g. compute
	Tags     circonus.Tags `json:"tags" toml:"tags" yaml:"tags"`             // service tags
	Disabled bool          `json:"disabled" toml:"disabled" yaml:"disabled"` // disable metric collection for this gcp service
}

// Filter defines any filtering criteria used for resources and metrics. Use labels, an expression, or neither.
type Filter struct {
	Labels     map[string]string `json:"labels" toml:"labels" yaml:"labels"`             // filter by labels
	Expression string            `json:"expression" toml:"expression" yaml:"expression"` // filter with a custom expression
}

// New creates a new collector instance.
func New(ctx context.Context, check *circonus.Check, cfgs []GCPCollector, interval time.Duration, logger zerolog.Logger) ([]Collector, error) {
	cl := collectorList()
	cc := []Collector{}

	if len(cfgs) == 0 {
		cfgs = []GCPCollector{
			{
				Name:     "compute",
				Disabled: false,
			},
		}
	}

	for _, cfg := range cfgs {
		cfg := cfg
		if cfg.Disabled {
			continue // entire service disabled
		}
		var c Collector
		var err error

		if initfn, known := cl[strings.ToLower(cfg.Name)]; known {
			c, err = initfn(ctx, check, &cfg, interval, logger)
		} else {
			err = errors.New("unrecognized aws service namespace")
		}

		if err != nil {
			logger.Warn().Err(err).Str("name", cfg.Name).Msg("skipping")
			continue
		}

		if c == nil {
			logger.Warn().Err(errors.New("init returned nil collector")).Str("name", cfg.Name).Msg("skipping")
			continue
		}

		cc = append(cc, c)
	}

	if len(cc) == 0 {
		return nil, errors.New("no collectors configured")
	}

	return cc, nil
}

type collectorInitFn func(context.Context, *circonus.Check, *GCPCollector, time.Duration, zerolog.Logger) (Collector, error)
type collectorInitList map[string]collectorInitFn

func collectorList() collectorInitList {
	return collectorInitList{
		"compute": newCompute,
	}
}

// ConfigExample generates configuration examples for collectors.
func ConfigExample() ([]GCPCollector, error) {
	var cc []GCPCollector
	cl := collectorList()
	for cn := range cl {
		c := GCPCollector{
			Name:     cn,
			Disabled: true,
		}
		switch cn {
		case "compute":
			c.Disabled = false
			cc = append(cc, c)
		default:
			return nil, errors.Errorf("unknown gcp service (%s)", cn)
		}
	}

	return cc, nil
}

type common struct {
	tsStart      time.Time
	tsEnd        time.Time
	ctx          context.Context
	disableTime  *time.Time
	check        *circonus.Check
	filter       Filter
	id           string
	disableCause string
	tags         circonus.Tags
	logger       zerolog.Logger
	interval     time.Duration
	enabled      bool
}

func newCommon(ctx context.Context, check *circonus.Check, cfg *GCPCollector, interval time.Duration, logger zerolog.Logger) common {
	return common{
		id:           cfg.Name,
		enabled:      true,
		disableCause: "",
		disableTime:  nil,
		check:        check,
		interval:     interval,
		filter:       cfg.Filter,
		ctx:          ctx,
		tags:         cfg.Tags,
		logger:       logger.With().Str("collector", cfg.Name).Logger(),
	}
}

func (c *common) ID() string {
	return c.id
}

func (c *common) done() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
		return false
	}
}
