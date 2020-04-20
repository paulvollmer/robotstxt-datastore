// Copyright 2020 Paul Vollmer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/paulvollmer/robotstxt"
	pb "github.com/paulvollmer/robotstxt-datastore/proto/robotstxt"
	"github.com/paulvollmer/robotstxt-datastore/server/ent"
)

func buildRobotstxtURL(host string) string {
	return cfg.DefaultRequestScheme + "://" + host + "/robots.txt"
}

// we do not use http.header.Get becasue we want to check case insensitive
func isHeaderContentTypeTextPlain(header http.Header) bool {
	for k, v := range header {
		if strings.ToLower(k) == "content-type" {
			if strings.Contains(strings.Join(v, " "), "text/plain") {
				return true
			}
		}
	}
	return false
}

// MapDBtoProtoResponse maps the database model to the protpbuf Robotstxt object
func MapDBtoProtoResponse(data *ent.Robotstxt) (*pb.Robotstxt, error) {
	createdAt, err := ptypes.TimestampProto(data.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := ptypes.TimestampProto(data.UpdatedAt)
	if err != nil {
		return nil, err
	}

	returnData := pb.Robotstxt{
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Host:         data.Host,
		Scheme:       data.Scheme,
		Statuscode:   data.Statuscode,
		ResponseUrl:  data.ResponseURL,
		ResponseTime: data.ResponseTime,
		Robotstxt:    data.Body,
	}

	if len(data.Body) == 0 {
		return &returnData, nil
	}

	robotstxtParsed, err := robotstxt.FromBytes(data.Body)
	if err != nil {
		return &returnData, err
	}

	returnData.Sitemaps = robotstxtParsed.Sitemaps

	rules := make([]*pb.Rule, 0)
	for k := range robotstxtParsed.Groups {
		paths := make([]string, 0)
		for i := 0; i < len(robotstxtParsed.Groups[k].Rules); i++ {
			trimmedPath := strings.TrimSpace(robotstxtParsed.Groups[k].Rules[i].Path)
			if trimmedPath != "" {
				paths = append(paths, trimmedPath)
			}
		}
		rules = append(rules, &pb.Rule{Agent: k, Paths: paths})
	}
	returnData.Rules = rules

	return &returnData, nil
}

func sentryCaptureException(err error, host, kind string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("host", host)
		scope.SetTag("kind", kind)
	})
	sentry.CaptureException(err)
}
