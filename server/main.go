// Copyright 2020 Paul Vollmer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/paulvollmer/robotstxt-datastore/proto/health"
	"github.com/paulvollmer/robotstxt-datastore/proto/robotstxt"
	"github.com/paulvollmer/robotstxt-datastore/server/ent"
	"github.com/paulvollmer/robotstxt-datastore/server/version"
	"google.golang.org/grpc"
)

var (
	dbClient *ent.Client
	cfg      Config
)

// Config store the environment variables to configure the server
type Config struct {
	ServerAddress        string `envconfig:"SERVER_ADDR"       default:":5000"`
	DatabaseHost         string `envconfig:"DATABASE_HOST"     default:"localhost"`
	DatabasePort         string `envconfig:"DATABASE_PORT"     default:"5432"`
	DatabaseUser         string `envconfig:"DATABASE_USER"     default:"postgres"`
	DatabasePassword     string `envconfig:"DATABASE_PASSWORD" default:"password"`
	DatabaseName         string `envconfig:"DATABASE_NAME"     default:"robotstxt"`
	DatabaseSSLMode      string `envconfig:"DATABASE_SSLMODE"  default:"disable"`
	DefaultRefreshAfter  int32  `envconfig:"REFRESH_AFTER"     default:"864000"`
	DefaultRequestScheme string `envconfig:"DEFAULT_REQUEST_SCHEME"  default:"https"`
	DefaultLimit         int    `envconfig:"DEFAULT_LIMIT"     default:"100"`
	UserAgent            string `envconfig:"USERAGENT"         default:"robotstxtbot"`
	SentryDSN            string `envconfig:"SENTRY_DSN"        default:""`
}

func main() {
	flagVersion := flag.Bool("version", false, "print the version and exit")
	flagEnv := flag.Bool("env", false, "print the environment variables and exit")
	flag.Parse()
	if *flagVersion {
		version.PrintInfo()
		os.Exit(0)
	}

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to configure: %v", err)
	}
	if *flagEnv {
		log.Printf("==> SERVER_ADDR            : %v\n", cfg.ServerAddress)
		log.Printf("==> DATABASE_HOST          : %v\n", cfg.DatabaseHost)
		log.Printf("==> DATABASE_PORT          : %v\n", cfg.DatabasePort)
		log.Printf("==> DATABASE_USER          : %v\n", cfg.DatabaseUser)
		log.Printf("==> DATABASE_PASSWORD      : %v\n", cfg.DatabasePassword)
		log.Printf("==> DATABASE_NAME          : %v\n", cfg.DatabaseName)
		log.Printf("==> DATABASE_SSLMODE       : %v\n", cfg.DatabaseSSLMode)
		log.Printf("==> REFRESH_AFTER          : %v\n", cfg.DefaultRefreshAfter)
		log.Printf("==> DEFAULT_REQUEST_SCHEME : %v\n", cfg.DefaultRequestScheme)
		log.Printf("==> DEFAULT_LIMIT          : %v\n", cfg.DefaultLimit)
		log.Printf("==> USERAGENT              : %v\n", cfg.UserAgent)
		log.Printf("==> SENTRY_DSN             : %v\n", cfg.SentryDSN)
		os.Exit(0)
	}

	// initialize sentry
	err = sentry.Init(sentry.ClientOptions{
		Dsn: cfg.SentryDSN,
	})
	if err != nil {
		log.Fatalf("error sentry init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	// initialize database connection
	driverName := "postgres"
	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName, cfg.DatabaseSSLMode)
	dbClient, err = ent.Open(driverName, dataSourceName)
	if err != nil {
		time.Sleep(5 * time.Second)
		log.Printf("retry 1 opening connection to postgres: %v", err)

		dbClient, err = ent.Open(driverName, dataSourceName)
		if err != nil {
			time.Sleep(10 * time.Second)
			log.Printf("retry 2 opening connection to postgres: %v", err)

			dbClient, err = ent.Open(driverName, dataSourceName)
			if err != nil {
				time.Sleep(20 * time.Second)
				log.Printf("retry 3 opening connection to postgres: %v", err)

				dbClient, err = ent.Open(driverName, dataSourceName)
				if err != nil {
					sentry.CaptureException(err)
					log.Fatalf("failed opening connection to postgres: %v", err)
				}
			}
		}
	}
	defer dbClient.Close()

	// run the auto migration tool.
	if err := dbClient.Schema.Create(context.Background()); err != nil {
		time.Sleep(5 * time.Second)
		log.Printf("retry 1 creating schema resources: %v", err)

		if err := dbClient.Schema.Create(context.Background()); err != nil {
			time.Sleep(10 * time.Second)
			log.Printf("retry 2 creating schema resources: %v", err)

			if err := dbClient.Schema.Create(context.Background()); err != nil {
				time.Sleep(20 * time.Second)
				log.Printf("retry 3 creating schema resources: %v", err)

				if err := dbClient.Schema.Create(context.Background()); err != nil {
					sentry.CaptureException(err)
					log.Fatalf("failed creating schema resources: %v", err)
				}
			}
		}
	}

	// register listener
	lis, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("failed to listen: %v", err)
	}
	// register Service
	grpcServer := grpc.NewServer()
	robotstxt.RegisterRobotstxtServiceServer(grpcServer, &serviceRobotstxt{})
	health.RegisterHealthServer(grpcServer, &serviceHealth{})

	// start the server
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			sentry.CaptureException(err)
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	log.Printf("server listening on %s\n", cfg.ServerAddress)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping grpc service...")
	grpcServer.GracefulStop()
}
