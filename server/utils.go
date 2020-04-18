// Copyright 2020 Paul Vollmer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/paulvollmer/robotstxt"
	pb "github.com/paulvollmer/robotstxt-datastore/proto/robotstxt"
	"github.com/paulvollmer/robotstxt-datastore/server/ent"
)

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

	robotstxtParsed, err := robotstxt.FromBytes(data.Body)
	if err != nil {
		return nil, err
	}

	groups := robotstxtParsed.Groups

	rules := make([]*pb.Rule, 0)
	for k := range groups {
		paths := make([]string, 0)
		for i := 0; i < len(groups[k].Rules); i++ {
			trimmedPath := strings.TrimSpace(groups[k].Rules[i].Path)
			if trimmedPath != "" {
				paths = append(paths, trimmedPath)
			}
		}
		rules = append(rules, &pb.Rule{Agent: k, Paths: paths})
	}

	return &pb.Robotstxt{
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Host:         data.Host,
		Scheme:       data.Scheme,
		Statuscode:   data.Statuscode,
		ResponseUrl:  data.ResponseURL,
		ResponseTime: data.ResponseTime,
		Robotstxt:    string(data.Body),
		Rules:        rules,
		Sitemaps:     robotstxtParsed.Sitemaps,
	}, nil
}
