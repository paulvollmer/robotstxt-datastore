// Copyright 2020 Paul Vollmer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/paulvollmer/robotstxt-datastore/server/ent"
	"github.com/stretchr/testify/assert"
)

func TestBuildRobotstxtURL(t *testing.T) {
	cfg.DefaultRequestScheme = "https"
	result := buildRobotstxtURL("text.com")
	assert.Equal(t, "https://text.com/robots.txt", result)
}

func TestIsHeaderContentTypeTextPlain(t *testing.T) {
	testData := []struct {
		in  http.Header
		out bool
	}{
		{
			in:  http.Header{"CONTENT-TYPE": []string{"text/plain"}},
			out: true,
		},
		{
			in:  http.Header{"Content-Type": []string{"text/plain"}},
			out: true,
		},
		{
			in:  http.Header{"Content-type": []string{"text/plain"}},
			out: true,
		},
		{
			in:  http.Header{"content-type": []string{"text/plain"}},
			out: true,
		},
		{
			in:  http.Header{"content-type": []string{"text/plain; charset=UTF-8"}},
			out: true,
		},
		{
			in:  http.Header{"test": []string{"this"}},
			out: false,
		},
		{
			in:  http.Header{"content-type": []string{"text/html"}},
			out: false,
		},
	}
	for index, tt := range testData {
		t.Run(strconv.Itoa(index), func(t *testing.T) {
			result := isHeaderContentTypeTextPlain(tt.in)
			assert.Equal(t, tt.out, result)
		})
	}
}

func TestMapDBtoProtoResponse(t *testing.T) {
	t.Run("body empty", func(t *testing.T) {
		data := ent.Robotstxt{
			Host:         "example.com",
			Scheme:       "https",
			Statuscode:   200,
			ResponseURL:  "https://example.com/robots.txt",
			ResponseTime: 100,
			Body:         []byte(""),
		}
		result, err := MapDBtoProtoResponse(&data)
		assert.Nil(t, err)
		assert.Equal(t, "example.com", result.Host)
		assert.Equal(t, "https", result.Scheme)
		assert.Equal(t, int32(200), result.Statuscode)
		assert.Equal(t, "https://example.com/robots.txt", result.ResponseUrl)
		assert.Equal(t, int64(100), result.ResponseTime)
		assert.Equal(t, []byte(""), result.Robotstxt)
		assert.Len(t, result.Rules, 0)
		assert.Len(t, result.Sitemaps, 0)
	})

	t.Run("body valid", func(t *testing.T) {
		body := []byte("Sitemap: https://example.com/sitemap1.xml\nSitemap: https://example.com/sitemap2.xml")
		data := ent.Robotstxt{
			Host:         "example.com",
			Scheme:       "https",
			Statuscode:   200,
			ResponseURL:  "https://example.com/robots.txt",
			ResponseTime: 100,
			Body:         body,
		}
		result, err := MapDBtoProtoResponse(&data)
		assert.Nil(t, err)
		assert.Equal(t, "example.com", result.Host)
		assert.Equal(t, "https", result.Scheme)
		assert.Equal(t, int32(200), result.Statuscode)
		assert.Equal(t, "https://example.com/robots.txt", result.ResponseUrl)
		assert.Equal(t, int64(100), result.ResponseTime)
		assert.Equal(t, body, result.Robotstxt)
		assert.Len(t, result.Rules, 0)
		assert.Len(t, result.Sitemaps, 2)
		assert.ElementsMatch(t, []string{"https://example.com/sitemap1.xml", "https://example.com/sitemap2.xml"}, result.Sitemaps)
	})

	t.Run("body invalid", func(t *testing.T) {
		body := "Disallow: /\nUser-agent: bot"
		data := ent.Robotstxt{
			Host:         "example.com",
			Scheme:       "https",
			Statuscode:   200,
			ResponseURL:  "https://example.com/robots.txt",
			ResponseTime: 100,
			Body:         []byte(body),
		}
		_, err := MapDBtoProtoResponse(&data)
		assert.NotNil(t, err)
	})
}
