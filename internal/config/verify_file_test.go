// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifyFile(t *testing.T) {
	tests := []struct {
		id          string
		file        string
		shouldFail  bool
		expectedErr string
	}{
		{"invalid - empty", "", true, "invalid file name (empty)"},
		{"invalid - missing", filepath.Join("testdata", "missing"), true, "no such file or directory"},
		{"invalid - not a file", filepath.Join("testdata", "not_a_file"), true, "not a regular file"},
		{"valid", filepath.Join("testdata", "test.file"), false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.id, func(t *testing.T) {
			t.Parallel()
			_, err := VerifyFile(tst.file)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				} else if tst.expectedErr != "" && !strings.Contains(err.Error(), tst.expectedErr) {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				}
			}
		})
	}
}
