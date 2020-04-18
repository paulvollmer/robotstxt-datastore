// Copyright 2020 Paul Vollmer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	pb "github.com/paulvollmer/robotstxt-datastore/proto/robotstxt"
	"github.com/stretchr/testify/assert"
)

func genRobotstxt(host, useragent string) string {
	return `# test robots.txt body
		User-agent: ` + useragent + `
		Disallow: /admin
		User-agent: *
		Disallow: /signup
		Sitemap: ` + host + `/sitemap.xml`
}

func TestCheckRobotstxt(t *testing.T) {
	cfg.DefaultRequestScheme = "http"

	t.Run("host empty", func(t *testing.T) {
		testSetup()

		var service serviceRobotstxt
		result, err := service.CheckRobotstxt(context.Background(), &pb.Check{Host: ""})
		assert.NotNil(t, err)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = host cannot be empty", err.Error())
		assert.Nil(t, result)

		testTeardown()
	})

	t.Run("host server connection refused", func(t *testing.T) {
		testSetup()

		var service serviceRobotstxt
		result, err := service.CheckRobotstxt(context.Background(), &pb.Check{Host: "localhost:1337"})
		assert.Nil(t, err)
		assert.Equal(t, "localhost:1337", result.Host)
		assert.Equal(t, "http", result.Scheme)
		assert.Equal(t, int32(0), result.Statuscode)
		//assert.NotEqual(t, int64(0), result.ResponseTime)
		assert.Equal(t, "", result.Robotstxt)
		assert.Len(t, result.Sitemaps, 0)
		assert.Len(t, result.Rules, 0)

		testTeardown()
	})

	t.Run("statuscode 404", func(t *testing.T) {
		testSetup()

		// create a mock http server without a robots.txt resource
		router := http.NewServeMux()
		router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
			// sleep a few milliseconds to set the response time to not zero
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(404)
		})
		testserver404 := httptest.NewServer(router)
		defer testserver404.Close()
		testserver404URL, err := url.Parse(testserver404.URL)
		if err != nil {
			panic(err)
		}

		var service serviceRobotstxt
		result, err := service.CheckRobotstxt(context.Background(), &pb.Check{Host: testserver404URL.Host})
		assert.Nil(t, err)
		assert.Equal(t, testserver404URL.Host, result.Host)
		//assert.Equal(t, "https", result.Scheme)
		assert.Equal(t, int32(404), result.Statuscode)
		assert.NotEqual(t, int64(0), result.ResponseTime)
		assert.Equal(t, "", result.Robotstxt)
		assert.Len(t, result.Sitemaps, 0)
		assert.Len(t, result.Rules, 0)

		testTeardown()
	})

	var testUrl string
	var useragent string
	// create a mock http server
	router := http.NewServeMux()
	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		// sleep a few milliseconds to set the response time to not zero
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, genRobotstxt(testUrl, useragent))
	})

	testserver := httptest.NewServer(router)
	defer testserver.Close()
	testUrl = testserver.URL
	testUrlParsed, err := url.Parse(testUrl)
	if err != nil {
		panic(err)
	}
	testHost := testUrlParsed.Host
	useragent = "testbot"

	resultEqual := func(t *testing.T, result *pb.Robotstxt) {
		assert.Equal(t, testHost, result.Host)
		assert.Equal(t, "http", result.Scheme)
		assert.Equal(t, int32(200), result.Statuscode)
		assert.NotEqual(t, int64(0), result.ResponseTime)
		assert.Equal(t, genRobotstxt(testUrl, useragent)+"\n", result.Robotstxt)
		assert.Len(t, result.Sitemaps, 1)
		assert.Equal(t, testUrl+"/sitemap.xml", result.Sitemaps[0])
		assert.Len(t, result.Rules, 2)
		assert.ElementsMatch(t, []*pb.Rule{
			{Agent: "*", Paths: []string{"/signup"}},
			{Agent: useragent, Paths: []string{"/admin"}},
		}, result.Rules)
	}

	t.Run("statuscode 200", func(t *testing.T) {
		testSetup()

		var service serviceRobotstxt
		result, err := service.CheckRobotstxt(context.Background(), &pb.Check{Host: testHost})
		assert.Nil(t, err)
		resultEqual(t, result)

		testTeardown()
	})

	t.Run("refresh", func(t *testing.T) {
		testSetup()

		var service serviceRobotstxt
		result, err := service.CheckRobotstxt(context.Background(), &pb.Check{Host: testHost})
		assert.Nil(t, err)
		resultEqual(t, result)

		useragent = "newtestbot"
		result2, err2 := service.CheckRobotstxt(context.Background(), &pb.Check{Host: testHost, Refresh: true})
		assert.Nil(t, err2)
		resultEqual(t, result2)

		testTeardown()
	})

	t.Run("refresh after 2", func(t *testing.T) {
		testSetup()

		var service serviceRobotstxt
		result, err := service.CheckRobotstxt(context.Background(), &pb.Check{Host: testHost})
		assert.Nil(t, err)
		resultEqual(t, result)

		time.Sleep(2 * time.Second)
		useragent = "newtestbot"
		result2, err2 := service.CheckRobotstxt(context.Background(), &pb.Check{Host: testHost, RefreshAfter: 1})
		assert.Nil(t, err2)
		resultEqual(t, result2)

		testTeardown()
	})
}

func TestListRobotstxts(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		testSetup()
		// mock multiple database entries
		for i := 0; i < 100; i++ {
			host := fmt.Sprintf("test-%v.com", i)
			mockBody := []byte(genRobotstxt(host, "testbot"))
			_, err := dbClient.Robotstxt.Create().
				SetHost(host).
				SetScheme("https").
				SetStatuscode(200).
				SetResponseTime(100).
				SetResponseURL("https://" + host + "/robots.txt").
				SetBody(mockBody).
				Save(context.Background())
			assert.Nil(t, err)
		}

		var service serviceRobotstxt
		result, err := service.ListRobotstxts(context.Background(), &pb.ListRobotstxtsRequest{Limit: 10, Offset: 0, SortBy: pb.ListRobotstxtsRequest_CREATED_AT})
		assert.Nil(t, err)
		assert.Equal(t, int64(100), result.TotalCount)
		assert.Len(t, result.Items, 10)
		for i := 0; i < 10; i++ {
			assert.Equal(t, fmt.Sprintf("test-%v.com", i), result.Items[i].Host)
			assert.Equal(t, "https", result.Items[i].Scheme)
			assert.Equal(t, int32(200), result.Items[i].Statuscode)
			assert.Equal(t, int64(100), result.Items[i].ResponseTime)
			assert.Equal(t, fmt.Sprintf("https://test-%v.com/robots.txt", i), result.Items[i].ResponseUrl)
			assert.Len(t, result.Items[i].Sitemaps, 1)
			assert.ElementsMatch(t, []string{fmt.Sprintf("test-%v.com/sitemap.xml", i)}, result.Items[i].Sitemaps)
			assert.Len(t, result.Items[i].Rules, 2)
			assert.ElementsMatch(t, []*pb.Rule{
				{Agent: "*", Paths: []string{"/signup"}},
				{Agent: "testbot", Paths: []string{"/admin"}},
			}, result.Items[i].Rules)
		}

		testTeardown()
	})

	t.Run("offset", func(t *testing.T) {
		t.Skip()
	})

	t.Run("search found", func(t *testing.T) {
		t.Skip()
	})

	t.Run("search not found", func(t *testing.T) {
		t.Skip()
	})
}

func TestReadRobotstxt(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		testSetup()
		// mock a database entry
		mockBody := genRobotstxt("test.com", "testbot")
		_, err := dbClient.Robotstxt.Create().
			SetHost("test.com").
			SetScheme("https").
			SetStatuscode(200).
			SetResponseTime(100).
			SetResponseURL("https://test.com/robots.txt").
			SetBody([]byte(mockBody)).
			Save(context.Background())
		assert.Nil(t, err)

		var service serviceRobotstxt
		in := &pb.GetRobotstxtRequest{Host: "test.com"}
		result, err := service.GetRobotstxt(context.Background(), in)
		assert.Nil(t, err)
		assert.Equal(t, "test.com", result.Host)
		assert.Equal(t, "https", result.Scheme)
		assert.Equal(t, int32(200), result.Statuscode)
		assert.Equal(t, int64(100), result.ResponseTime)
		assert.Equal(t, "https://test.com/robots.txt", result.ResponseUrl)
		assert.Equal(t, mockBody, result.Robotstxt)
		assert.Len(t, result.Sitemaps, 1)
		assert.ElementsMatch(t, []string{"test.com/sitemap.xml"}, result.Sitemaps)
		assert.Len(t, result.Rules, 2)
		assert.ElementsMatch(t, []*pb.Rule{
			{Agent: "*", Paths: []string{"/signup"}},
			{Agent: "testbot", Paths: []string{"/admin"}},
		}, result.Rules)

		testTeardown()
	})

	t.Run("not found", func(t *testing.T) {
		testSetup()

		var service serviceRobotstxt
		in := &pb.GetRobotstxtRequest{Host: "notfound.com"}
		result, err := service.GetRobotstxt(context.Background(), in)
		assert.NotNil(t, err)
		assert.Equal(t, "rpc error: code = NotFound desc = host not found", err.Error())
		assert.Nil(t, result)

		testTeardown()
	})
}
