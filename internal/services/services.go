// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package services

type Service interface {
	Enabled() bool
	Scan() error
	Start() error
}
