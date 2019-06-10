// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package config

// func TestValidateAPIOptions(t *testing.T) {
// 	t.Log("Testing validateAPIOptions")
//
// 	t.Log("No key/app/url")
// 	{
// 		expectedError := errors.New("API key is required")
// 		err := validateAPIOptions()
// 		if err == nil {
// 			t.Fatal("Expected error")
// 		}
// 		if err.Error() != expectedError.Error() {
// 			t.Errorf("Expected (%s) got (%s)", expectedError, err)
// 		}
// 	}
//
// 	t.Log("No app")
// 	{
// 		viper.Set(KeyAPITokenKey, "foo")
// 		expectedError := errors.New("API app is required")
// 		err := validateAPIOptions()
// 		if err == nil {
// 			t.Fatal("Expected error")
// 		}
// 		if err.Error() != expectedError.Error() {
// 			t.Errorf("Expected (%s) got (%s)", expectedError, err)
// 		}
// 	}
//
// 	t.Log("No url")
// 	{
// 		viper.Set(KeyAPITokenKey, "foo")
// 		viper.Set(KeyAPITokenApp, "foo")
// 		expectedError := errors.New("API URL is required")
// 		err := validateAPIOptions()
// 		if err == nil {
// 			t.Fatal("Expected error")
// 		}
// 		if err.Error() != expectedError.Error() {
// 			t.Errorf("Expected (%s) got (%s)", expectedError, err)
// 		}
// 	}
//
// 	t.Log("Invalid url (foo)")
// 	{
// 		viper.Set(KeyAPITokenKey, "foo")
// 		viper.Set(KeyAPITokenApp, "foo")
// 		viper.Set(KeyAPIURL, "foo")
// 		expectedError := errors.New("invalid API URL (foo)")
// 		err := validateAPIOptions()
// 		if err == nil {
// 			t.Fatal("Expected error")
// 		}
// 		if err.Error() != expectedError.Error() {
// 			t.Errorf("Expected (%s) got (%s)", expectedError, err)
// 		}
// 	}
//
// 	t.Log("Invalid url (foo_bar://herp/derp)")
// 	{
// 		viper.Set(KeyAPITokenKey, "foo")
// 		viper.Set(KeyAPITokenApp, "foo")
// 		viper.Set(KeyAPIURL, "foo_bar://herp/derp")
// 		expectedError := errors.New("invalid API URL: parse foo_bar://herp/derp: first path segment in URL cannot contain colon")
// 		err := validateAPIOptions()
// 		if err == nil {
// 			t.Fatal("Expected error")
// 		}
// 		if err.Error() != expectedError.Error() {
// 			t.Errorf("Expected (%s) got (%s)", expectedError, err)
// 		}
// 	}
//
// 	t.Log("Valid options")
// 	{
// 		viper.Set(KeyAPITokenKey, "foo")
// 		viper.Set(KeyAPITokenApp, "foo")
// 		viper.Set(KeyAPIURL, "http://foo.com/bar")
// 		err := validateAPIOptions()
// 		if err != nil {
// 			t.Fatalf("Expected NO error, got (%s)", err)
// 		}
// 	}
//
// }
