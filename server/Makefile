# Copyright 2020 Paul Vollmer. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

all: fmt lint test build

VERSION=0.1.0

LDFLAGS := --ldflags "-s -w \
-X github.com/paulvollmer/robotstxt-datastore/version.Version=${VERSION} \
-X github.com/paulvollmer/robotstxt-datastore/version.GitCommit=$(shell git rev-parse --short HEAD 2> /dev/null || true) \
-X github.com/paulvollmer/robotstxt-datastore/version.BuildDate=$(shell date +%Y.%m.%d.%H%M%S) \
-X github.com/paulvollmer/robotstxt-datastore/version.BuildOS=$(shell go env GOHOSTOS) \
-X github.com/paulvollmer/robotstxt-datastore/version.BuildArch=$(shell go env GOHOSTARCH)"

fmt:
	@go fmt ./...

lint:
	@golangci-lint run

ent:
	@entc generate ./ent/schema

test: clean
	@go test -cover ./...

test-cov: clean
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

build:
	@go build ${LDFLAGS} -o robotstxt-datastore

build-darwin:
	@GOOS=darwin GOARCH=386   go build ${LDFLAGS} -o robotstxt-datastore-darwin-386
	@GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o robotstxt-datastore-darwin-amd64

build-linux:
	@GOOS=linux GOARCH=386   go build ${LDFLAGS} --tags production -o robotstxt-datastore-linux-386
	@GOOS=linux GOARCH=amd64 go build ${LDFLAGS} --tags production -o robotstxt-datastore-linux-amd64

clean:
	@rm -f robotstxt-datastore robotstxt-datastore-darwin-* robotstxt-datastore-linux-*

.PHONY: fmt lint ent test test-cov clean build build-darwin build-linux
