# Copyright 2020 Paul Vollmer. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# build stage
FROM golang AS builder
ADD . /app
WORKDIR /app
ENV GO111MODULE=on
RUN cd server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /server .

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /server ./
RUN chmod +x ./server
ENTRYPOINT ["./server"]
EXPOSE 5000
