# Copyright 2020 Paul Vollmer. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# build stage
FROM golang:1.14.2 AS builder
ADD . /workspace
WORKDIR /workspace
ENV GO111MODULE=on
RUN cd server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /workspace/server/robotstxt-datastore ./robotstxt-datastore
RUN chmod +x ./robotstxt-datastore
RUN ./robotstxt-datastore -version
ENTRYPOINT ["./robotstxt-datastore"]
EXPOSE 5000
