// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package circonus

import (
	"bytes"
)

// ReportError to Circonus check.
func (c *Check) ReportError(err error) {
	var buf bytes.Buffer

	if e := c.WriteMetricSample(&buf, c.errorMetricName, MetricTypeString, err.Error(), nil); err != nil {
		c.logger.Error().Err(e).Msg("writing error metric sample")
		return
	}

	if e := c.SubmitMetrics(&buf); e != nil {
		c.logger.Error().Err(e).Msg("submitting error metric sample")
	}
}
