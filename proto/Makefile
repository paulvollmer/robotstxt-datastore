# Copyright 2020 Paul Vollmer. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

all: lint build

PATH_ALL := --proto_path . --proto_path /usr/local/include
PROTOC_FLAGS := ${PATH_ALL} --go_out=plugins=grpc:.

build: docs
	@protoc ${PROTOC_FLAGS} robotstxt/*.proto
	@protoc ${PROTOC_FLAGS} health/*.proto

lint:
	@protoc ${PATH_ALL} --lint_out=. robotstxt/*.proto health/*.proto

clean: docs-clean
	@rm -f health/health.pb.go
	@rm -f robotstxt/robotstxt.pb.go

docs:
	@rm -rf ../docs
	@mkdir -p ../docs
	@protoc ${PATH_ALL} --doc_out=../docs --doc_opt=html,index.html robotstxt/*.proto health/*.proto

.PHONY: all build lint clean docs
