// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/circonus-cloud-agent/internal/config"
	"github.com/circonus-labs/circonus-cloud-agent/internal/config/defaults"
	"github.com/circonus-labs/circonus-cloud-agent/internal/release"
	"github.com/spf13/viper"
)

//
// Common configuration settings
//

func init() {
	//
	// General settings
	//
	{
		const (
			key         = config.KeyDebug
			longOpt     = "debug"
			shortOpt    = "d"
			envVar      = release.ENVPREFIX + "_DEBUG"
			description = "Enable debug messages"
		)

		RootCmd.PersistentFlags().BoolP(longOpt, shortOpt, defaults.Debug, envDescription(description, envVar))
		if err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt)); err != nil {
			bindFlagError(longOpt, err)
		}
		if err := viper.BindEnv(key, envVar); err != nil {
			bindEnvError(envVar, err)
		}
		viper.SetDefault(key, defaults.Debug)
	}

	{
		const (
			key         = config.KeyLogLevel
			longOpt     = "log-level"
			envVar      = release.ENVPREFIX + "_LOG_LEVEL"
			description = "Log level [(panic|fatal|error|warn|info|debug|disabled)]"
		)

		RootCmd.PersistentFlags().String(longOpt, defaults.LogLevel, envDescription(description, envVar))
		if err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt)); err != nil {
			bindFlagError(longOpt, err)
		}
		if err := viper.BindEnv(key, envVar); err != nil {
			bindEnvError(envVar, err)
		}
		viper.SetDefault(key, defaults.LogLevel)
	}

	{
		const (
			key         = config.KeyLogPretty
			longOpt     = "log-pretty"
			envVar      = release.ENVPREFIX + "_LOG_PRETTY"
			description = "Output formatted/colored log lines [ignored on windows]"
		)

		RootCmd.PersistentFlags().Bool(longOpt, defaults.LogPretty, envDescription(description, envVar))
		if err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt)); err != nil {
			bindFlagError(longOpt, err)
		}
		if err := viper.BindEnv(key, envVar); err != nil {
			bindEnvError(envVar, err)
		}
		viper.SetDefault(key, defaults.LogPretty)
	}

	// {
	// 	const (
	// 		key         = config.KeyPipeSubmits
	// 		longOpt     = "pipe-submits"
	// 		envVar      = release.ENVPREFIX + "_PIPE_SUBMITS"
	// 		description = "Pipe metric submissions to Circonus (experimental)"
	// 	)

	// 	RootCmd.PersistentFlags().Bool(longOpt, false, envDescription(description, envVar))
	// 	viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
	// 	viper.BindEnv(key, envVar)
	// 	viper.SetDefault(key, false)
	// }
}
