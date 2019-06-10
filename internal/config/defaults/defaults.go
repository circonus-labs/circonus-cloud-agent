// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package defaults

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/circonus-labs/circonus-cloud-agent/internal/release"
)

const (
	// APIURL for circonus
	APIURL = "https://api.circonus.com/v2/"

	// APIApp defines the api app name associated with the api token key
	APIApp = release.NAME

	// APIDebug is false by default
	APIDebug = false

	// Debug is false by default
	Debug = false

	// MetricNameSeparator defines character used to delimit metric name parts
	MetricNameSeparator = "`"

	// LogLevel set to info by default
	LogLevel = "info"

	// LogPretty colored/formatted output to stderr
	LogPretty = false
)

var (
	// BasePath is the "base" directory
	//
	// expected installation structure:
	// base        (e.g. /opt/circonus/cloud-agent)
	//   /bin      (e.g. /opt/circonus/cloud-agent/bin)
	//   /etc      (e.g. /opt/circonus/cloud-agent/etc)
	//   /sbin     (e.g. /opt/circonus/cloud-agent/sbin)
	BasePath = ""

	// ConfigFile defines the default configuration file name
	ConfigFile = ""

	// EtcPath returns the default etc directory within base directory
	EtcPath = "" // (e.g. /opt/circonus/cloud-agent/etc)
)

func init() {
	var exePath string
	var resolvedExePath string
	var err error

	exePath, err = os.Executable()
	if err == nil {
		resolvedExePath, err = filepath.EvalSymlinks(exePath)
		if err == nil {
			BasePath = filepath.Clean(filepath.Join(filepath.Dir(resolvedExePath), "..")) // e.g. /opt/circonus/cloud-agent
		}
	}

	if err != nil {
		fmt.Printf("Unable to determine path to binary %v\n", err)
		os.Exit(1)
	}

	EtcPath = filepath.Join(BasePath, "etc")
	ConfigFile = filepath.Join(EtcPath, release.NAME+".yaml")
}
