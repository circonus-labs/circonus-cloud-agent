// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/circonus-cloud-agent/internal/services/azureservice"
	"github.com/spf13/viper"
)

//
// Azure configuration settings
//

func init() {
	{
		const (
			key         = azureservice.KeyEnabled
			longOpt     = "enable-azure"
			description = "Enable Azure metric collection client"
		)

		RootCmd.PersistentFlags().Bool(longOpt, azureservice.DefaultEnabled, description)
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.SetDefault(key, azureservice.DefaultEnabled)
	}

	{
		const (
			key         = azureservice.KeyConfDir
			longOpt     = "azure-conf-dir"
			description = "Azure configuration directory"
		)

		RootCmd.PersistentFlags().String(longOpt, azureservice.DefaultConfDir, description)
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.SetDefault(key, azureservice.DefaultConfDir)
	}
	{
		const (
			key         = azureservice.KeyConfExample
			longOpt     = "azure-example-conf"
			description = "Show Azure config (json|toml|yaml) and exit"
		)

		RootCmd.PersistentFlags().String(longOpt, "", description)
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
	}
}
