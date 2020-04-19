// Copyright 2020 Paul Vollmer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errorHostCannotBeEmpty = status.Error(codes.InvalidArgument, "host cannot be empty")
	errorHostNotFound      = status.Error(codes.NotFound, "host not found")
)
