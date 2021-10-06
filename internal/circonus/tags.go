// Copyright © 2019 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package circonus

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// Tag defines an individual tag.
type Tag struct {
	Category string
	Value    string
}

// Tags defines a list of tags.
type Tags []Tag

// MetricNameWithStreamTags will encode tags as stream tags into supplied metric name.
// Note: if metric name already has stream tags it is assumed the metric name and
// embedded stream tags are being managed manually and calling this method will nave no effect.
func (c *Check) MetricNameWithStreamTags(metric string, tags Tags) string {
	if len(tags) == 0 {
		return metric
	}

	if strings.Contains(metric, "|ST[") {
		return metric
	}

	taglist := c.EncodeMetricStreamTags(tags)
	if taglist != "" {
		return metric + "|ST[" + taglist + "]"
	}

	return metric
}

// EncodeMetricStreamTags encodes Tags into a string suitable for use
// with stream tags. Tags directly embedded into metric names using the
// `metric_name|ST[<tags>]` syntax. Note: legacy tags support bare
// values, stream tags require a cateogry and a value. Additionally,
// all spaces are removed from stream tag categories and values.
func (c *Check) EncodeMetricStreamTags(tags Tags) string {
	if len(tags) == 0 {
		return ""
	}

	tmpTags := c.EncodeMetricTags(tags)
	if len(tmpTags) == 0 {
		return ""
	}

	tagList := make([]string, len(tmpTags))
	for i, tag := range tmpTags {
		if i >= MaxTags {
			c.logger.Warn().Int("num", len(tags)).Int("max", MaxTags).Interface("tags", tags).Msg("ignoring tags over max")
			break
		}

		tagParts := strings.SplitN(tag, ":", 2)
		if len(tagParts) != 2 {
			c.logger.Warn().Str("tag", tag).Msg("stream tags must have a category and value, ignoring tag")
			continue // invalid tag, skip it
		}
		tc := tagParts[0]
		tv := tagParts[1]

		encodeFmt := `b"%s"`
		encodedSig := `b"` // has cat or val been previously (or manually) base64 encoded and formatted
		if !strings.HasPrefix(tc, encodedSig) {
			tc = fmt.Sprintf(encodeFmt, base64.StdEncoding.EncodeToString([]byte(strings.Map(removeSpaces, tc))))
		}
		if !strings.HasPrefix(tv, encodedSig) {
			tv = fmt.Sprintf(encodeFmt, base64.StdEncoding.EncodeToString([]byte(strings.Map(removeSpaces, tv))))
		}

		tagList[i] = tc + ":" + tv
	}

	return strings.Join(tagList, ",")
}

// EncodeMetricTags encodes Tags into an array of strings. The format
// check_bundle.metircs.metric.tags needs. This helper is intended to work
// with legacy check bundle metrics.
func (c *Check) EncodeMetricTags(tags Tags) []string {
	if len(tags) == 0 {
		return []string{}
	}

	uniqueTags := make(map[string]bool)
	for i, t := range tags {
		if i >= MaxTags {
			c.logger.Warn().Int("num", len(tags)).Int("max", MaxTags).Interface("tags", tags).Msg("max tags reached, ignoring remainder")
			break
		}

		// prep components, trim leading/trailing white space and lowercase
		tc := strings.ToLower(strings.TrimSpace(t.Category))
		tv := strings.ToLower(strings.TrimSpace(t.Value))
		if tc == "" && tv == "" {
			continue // empty tag (invalid), skip it
		}

		var sb strings.Builder
		if tc != "" {
			sb.WriteString(tc)
		}
		if tc != "" && tv != "" {
			sb.WriteString(":")
		}
		if tv != "" {
			sb.WriteString(tv)
		}
		uniqueTags[sb.String()] = true
	}

	// create list of unique tags and sort
	tagList := make([]string, len(uniqueTags))
	idx := 0
	for t := range uniqueTags {
		tagList[idx] = t
		idx++
	}
	sort.Strings(tagList)

	return tagList
}

func removeSpaces(r rune) rune {
	if unicode.IsSpace(r) {
		return -1
	}
	return r
}
