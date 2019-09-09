// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/circonus-cloud-agent/internal/services/awsservice"
	"github.com/spf13/viper"
)

//
// AWS configuration settings
//

func init() {
	{
		const (
			key         = awsservice.KeyEnabled
			longOpt     = "enable-aws"
			description = "Enable AWS metric collection client"
		)

		RootCmd.PersistentFlags().Bool(longOpt, awsservice.DefaultEnabled, description)
		if err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt)); err != nil {
			bindFlagError(longOpt, err)
		}
		viper.SetDefault(key, awsservice.DefaultEnabled)
	}

	{
		const (
			key         = awsservice.KeyConfDir
			longOpt     = "aws-conf-dir"
			description = "AWS configuration directory"
		)

		RootCmd.PersistentFlags().String(longOpt, awsservice.DefaultConfDir, description)
		if err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt)); err != nil {
			bindFlagError(longOpt, err)
		}
		viper.SetDefault(key, awsservice.DefaultConfDir)
	}
	{
		const (
			key         = awsservice.KeyConfExample
			longOpt     = "aws-example-conf"
			description = "Show AWS config (json|toml|yaml) and exit"
		)

		RootCmd.PersistentFlags().String(longOpt, "", description)
		if err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt)); err != nil {
			bindFlagError(longOpt, err)
		}
	}

}
