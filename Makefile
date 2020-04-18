# Copyright 2020 Paul Vollmer. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

VERSION := $(shell head -1 VERSION)

docker-up:
	@docker-compose up --build

docker-down:
	@docker-compose down

docker-clean:
	@rm -rf postgres-data

docker-build:
	@docker build -t paulvollmer/robotstxt-datastore:v${VERSION} .

docker-push: docker-build
	@docker push paulvollmer/robotstxt-datastore:v${VERSION}
	
.PHONY: docker-up docker-down docker-clean docker-build docker-push
