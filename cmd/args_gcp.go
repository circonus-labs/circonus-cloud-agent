// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/circonus-cloud-agent/internal/services/gcpservice"
	"github.com/spf13/viper"
)

//
// GCP configuration settings
//

func init() {
	{
		const (
			key         = gcpservice.KeyEnabled
			longOpt     = "enable-gcp"
			description = "Enable GCP metric collection client"
		)

		RootCmd.PersistentFlags().Bool(longOpt, gcpservice.DefaultEnabled, description)
		if err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt)); err != nil {
			bindFlagError(longOpt, err)
		}
		viper.SetDefault(key, gcpservice.DefaultEnabled)
	}

	{
		const (
			key         = gcpservice.KeyConfDir
			longOpt     = "gcp-conf-dir"
			description = "GCP configuration directory"
		)

		RootCmd.PersistentFlags().String(longOpt, gcpservice.DefaultConfDir, description)
		if err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt)); err != nil {
			bindFlagError(longOpt, err)
		}
		viper.SetDefault(key, gcpservice.DefaultConfDir)
	}
	{
		const (
			key         = gcpservice.KeyConfExample
			longOpt     = "gcp-example-conf"
			description = "Show GCP config (json|toml|yaml) and exit"
		)

		RootCmd.PersistentFlags().String(longOpt, "", description)
		if err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt)); err != nil {
			bindFlagError(longOpt, err)
		}
	}

}
