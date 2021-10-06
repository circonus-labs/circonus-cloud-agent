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
		expectedErr string
		shouldFail  bool
	}{
		{id: "invalid - empty", file: "", shouldFail: true, expectedErr: "invalid file name (empty)"},
		{id: "invalid - missing", file: filepath.Join("testdata", "missing"), shouldFail: true, expectedErr: "no such file or directory"},
		{id: "invalid - not a file", file: filepath.Join("testdata", "not_a_file"), shouldFail: true, expectedErr: "not a regular file"},
		{id: "valid", file: filepath.Join("testdata", "test.file"), shouldFail: false, expectedErr: ""},
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
