# Copyright 2020 Paul Vollmer. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

VERSION := $(shell head -1 VERSION)
DOCKER_REGISTRY := 192.168.0.39:9000

docker-up:
	@docker-compose up --build

docker-down:
	@docker-compose down

docker-clean:
	@rm -rf postgres-data

docker-build:
	@docker build -t ${DOCKER_REGISTRY}/robotstxt-datastore:v${VERSION} .

docker-push:
	@docker push ${DOCKER_REGISTRY}/robotstxt-datastore:v${VERSION}

.PHONY: docker-up docker-down docker-clean docker-build docker-push
