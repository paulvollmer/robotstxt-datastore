// Copyright 2020 Paul Vollmer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paulvollmer/robotstxt-datastore/server/ent"
)

func testSetup() {
	var err error

	// init database client
	dbClient, err = ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		panic(err)
	}
	// Run the auto migration tool.
	if err := dbClient.Schema.Create(context.Background()); err != nil {
		panic(err)
	}
}

func testTeardown() {
	dbClient.Close()
}
