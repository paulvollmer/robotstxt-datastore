// Copyright 2020 Paul Vollmer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	pb "github.com/paulvollmer/robotstxt-datastore/proto/robotstxt"
	"github.com/paulvollmer/robotstxt-datastore/server/ent"
	"github.com/paulvollmer/robotstxt-datastore/server/ent/robotstxt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serviceRobotstxt struct{}

func (s *serviceRobotstxt) CheckRobotstxt(ctx context.Context, in *pb.Check) (*pb.Robotstxt, error) {
	log.Println("CheckRobotstxt:", in)

	// if refreshAfter is zero, set to default value
	refreshAfter := in.GetRefreshAfter()
	if refreshAfter == 0 {
		refreshAfter = cfg.DefaultRefreshAfter
	}
	refresh := in.Refresh

	// is the host not empty
	host := in.GetHost()
	if host == "" {
		return nil, errorHostCannotBeEmpty
	}

	// exist the database entry?
	foundDB, err := dbClient.Robotstxt.Query().Where(robotstxt.Host(host)).All(ctx)
	if err != nil {
		sentryCaptureException(err, host, "db query")
		return nil, status.Error(codes.Internal, err.Error())
	}
	if len(foundDB) == 1 {
		found, err := MapDBtoProtoResponse(foundDB[0])
		if err != nil {
			sentryCaptureException(err, host, "map db")
			return nil, status.Error(codes.Internal, err.Error())
		}
		log.Println("--> Found at database, sitemaps total:", len(found.Sitemaps))

		after := foundDB[0].UpdatedAt.Add(time.Second * time.Duration(refreshAfter))
		if time.Now().After(after) {
			log.Println("--> Outdated robots.txt, Reload resource...")
			refresh = true
		} else {
			if !refresh {
				return found, nil
			}
		}
	}

	// fetch the robots.txt data
	targetUri := buildRobotstxtURL(host)
	robotstxtData := []byte{}
	var client http.Client

	// TODO: add env var to enable insecure requests
	//client.Transport =  &http.Transport{
	//	TLSClientConfig: &tls.Config{
	//		InsecureSkipVerify: true,
	//	},
	//}

	req, err := http.NewRequest("GET", targetUri, nil)
	if err != nil {
		sentryCaptureException(err, host, "req get")
		return nil, status.Error(codes.Internal, err.Error())
	}
	req.Header.Set("User-Agent", cfg.UserAgent)
	start := time.Now()
	resp, err := client.Do(req)
	responseTime := time.Since(start)
	if err != nil {
		// if a client error occur, create a database entry and set the StatusCode to zero
		data, err := dbClient.Robotstxt.Create().
			SetHost(host).
			SetScheme(cfg.DefaultRequestScheme).
			SetStatuscode(0).
			SetResponseTime(responseTime.Milliseconds()).
			SetResponseURL("").
			SetBody([]byte("")).
			Save(ctx)
		if err != nil {
			sentryCaptureException(err, host, "db create")
			return nil, status.Error(codes.Internal, err.Error())
		}
		res, err := MapDBtoProtoResponse(data)
		return res, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {

		if isHeaderContentTypeTextPlain(resp.Header) == false {
			log.Println("--> Content-Type is not text/plain")
			data, err := dbClient.Robotstxt.Create().
				SetHost(host).
				SetScheme(resp.Request.URL.Scheme).
				SetStatuscode(int32(resp.StatusCode)).
				SetResponseTime(responseTime.Milliseconds()).
				SetResponseURL(resp.Request.URL.String()).
				SetBody([]byte("")).
				Save(ctx)
			if err != nil {
				sentryCaptureException(err, host, "db create")
				return nil, status.Error(codes.Internal, err.Error())
			}
			res, err := MapDBtoProtoResponse(data)
			return res, err
		}

		robotstxtData, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			sentryCaptureException(err, host, "req read body")
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	// does not exist at the database
	if len(foundDB) == 0 {
		log.Println("--> Create database entry...")
		data, err := dbClient.Robotstxt.Create().
			SetHost(host).
			SetScheme(resp.Request.URL.Scheme).
			SetStatuscode(int32(resp.StatusCode)).
			SetResponseTime(responseTime.Milliseconds()).
			SetResponseURL(resp.Request.URL.String()).
			SetBody(robotstxtData).
			Save(ctx)
		if err != nil {
			sentryCaptureException(err, host, "db create")
			return nil, status.Error(codes.Internal, err.Error())
		}
		res, err := MapDBtoProtoResponse(data)
		return res, err
	}

	// refresh the database
	if refresh && len(foundDB) == 1 {
		log.Println("--> Update database entry...")
		data, err := foundDB[0].Update().
			SetScheme(resp.Request.URL.Scheme).
			SetStatuscode(int32(resp.StatusCode)).
			SetResponseTime(responseTime.Milliseconds()).
			SetResponseURL(resp.Request.URL.String()).
			SetBody(robotstxtData).
			Save(ctx)
		if err != nil {
			sentryCaptureException(err, host, "db update")
			return nil, status.Error(codes.Internal, err.Error())
		}
		res, err := MapDBtoProtoResponse(data)
		return res, err
	}

	return nil, nil
}

// ListRobotstxts return a list of robots.txt sources.
func (s *serviceRobotstxt) ListRobotstxts(ctx context.Context, in *pb.ListRobotstxtsRequest) (*pb.ListRobotstxtsResponse, error) {
	log.Println("ListRobotstxts:", in)

	limit := int(in.Limit)
	// set the default limit
	if limit == 0 {
		limit = cfg.DefaultLimit
	}
	offset := int(in.Offset)

	// get the data
	count, err := dbClient.Robotstxt.Query().Count(ctx)
	if err != nil {
		sentry.CaptureException(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	dbQuery := dbClient.Robotstxt.Query().Limit(limit).Offset(offset)

	search := in.GetSearch()
	if search != "" {
		dbQuery.Where(robotstxt.HostContains(search))
	}

	sort := ""
	switch in.GetSortBy() {
	case pb.ListRobotstxtsRequest_HOST:
		sort = robotstxt.FieldHost
	case pb.ListRobotstxtsRequest_CREATED_AT:
		sort = robotstxt.FieldCreatedAt
	case pb.ListRobotstxtsRequest_UPDATED_AT:
		sort = robotstxt.FieldUpdatedAt
	default:
		sort = robotstxt.FieldHost
	}

	switch in.GetOrderBy() {
	case pb.ListRobotstxtsRequest_ASC:
		dbQuery.Order(ent.Asc(sort))
	case pb.ListRobotstxtsRequest_DESC:
		dbQuery.Order(ent.Desc(sort))
	default:
		dbQuery.Order(ent.Asc(sort))
	}

	data, err := dbQuery.All(ctx)
	if err != nil {
		sentry.CaptureException(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// create the response data and return
	reply := &pb.ListRobotstxtsResponse{}
	reply.TotalCount = int64(count)
	tmpItems := make([]*pb.Robotstxt, len(data))
	for i := 0; i < len(data); i++ {
		d, err := MapDBtoProtoResponse(data[i])
		if err != nil {
			sentry.CaptureException(err)
			return nil, err
		}
		tmpItems[i] = d
	}
	reply.Items = tmpItems
	return reply, nil
}

// GetRobotstxt return the robots.txt response object for the given host.
func (s *serviceRobotstxt) GetRobotstxt(ctx context.Context, in *pb.GetRobotstxtRequest) (*pb.Robotstxt, error) {
	host := in.GetHost()
	log.Printf("GetRobotstxt Host: %q\n", host)

	if host == "" {
		return nil, errorHostCannotBeEmpty
	}

	data, err := dbClient.Robotstxt.Query().Where(robotstxt.Host(host)).All(ctx)
	if err != nil {
		sentryCaptureException(err, host, "db query")
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(data) == 0 {
		return nil, errorHostNotFound
	}

	res, err := MapDBtoProtoResponse(data[0])
	if err != nil {
		sentryCaptureException(err, host, "map db")
		return nil, err
	}
	return res, nil
}
