// Copyright 2020 Paul Vollmer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
	"time"

	"github.com/paulvollmer/robotstxt-datastore/server/ent"
	"github.com/stretchr/testify/assert"
)

func TestMapDBtoProtoResponse(t *testing.T) {
	body := `
		Sitemap: https://example.com/sitemap1.xml
		Sitemap: https://example.com/sitemap2.xml
	`
	data := ent.Robotstxt{
		ID:           0,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
		Host:         "example.com",
		Scheme:       "https",
		Body:         []byte(body),
		Statuscode:   200,
		ResponseTime: 100,
	}
	result, err := MapDBtoProtoResponse(&data)
	assert.Nil(t, err)
	assert.Equal(t, "example.com", result.Host)
	assert.Equal(t, "https", result.Scheme)
	assert.Equal(t, body, result.Robotstxt)
	assert.ElementsMatch(t, []string{"https://example.com/sitemap1.xml", "https://example.com/sitemap2.xml"}, result.Sitemaps)
	//result.Rules
	assert.Equal(t, int32(200), result.Statuscode)
	assert.Equal(t, int64(100), result.ResponseTime)
}
