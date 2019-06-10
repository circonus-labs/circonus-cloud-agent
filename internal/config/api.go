// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package config

// func validateAPIOptions() error {
// 	apiKey := viper.GetString(KeyAPITokenKey)
// 	apiApp := viper.GetString(KeyAPITokenApp)
// 	apiURL := viper.GetString(KeyAPIURL)
// 	apiCAFile := viper.GetString(KeyAPICAFile)
//
// 	// API is required for reverse and/or statsd
//
// 	if apiKey == "" {
// 		return errors.New("API key is required")
// 	}
//
// 	if apiApp == "" {
// 		return errors.New("API app is required")
// 	}
//
// 	if apiURL == "" {
// 		return errors.New("API URL is required")
// 	}
//
// 	if apiURL != defaults.APIURL {
// 		parsedURL, err := url.Parse(apiURL)
// 		if err != nil {
// 			return errors.Wrap(err, "invalid API URL")
// 		}
// 		if parsedURL.Scheme == "" || parsedURL.Host == "" || parsedURL.Path == "" {
// 			return errors.Errorf("invalid API URL (%s)", apiURL)
// 		}
// 	}
//
// 	// NOTE the api ca file doesn't come from the cosi config
// 	if apiCAFile != "" {
// 		f, err := verifyFile(apiCAFile)
// 		if err != nil {
// 			return err
// 		}
// 		viper.Set(KeyAPICAFile, f)
// 	}
//
// 	viper.Set(KeyAPITokenKey, apiKey)
// 	viper.Set(KeyAPITokenApp, apiApp)
// 	viper.Set(KeyAPIURL, apiURL)
//
// 	return nil
// }
