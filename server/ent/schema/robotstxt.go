// Copyright 2020 Paul Vollmer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

// Robotstxt holds the schema definition for the Robotstxt entity.
type Robotstxt struct {
	ent.Schema
}

// Fields of the Robotstxt.
func (Robotstxt) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.String("host").Unique(),
		field.String("scheme"),
		field.String("response_url"),
		field.Bytes("body"),
		field.Int32("statuscode"),
		field.Int64("response_time"),
	}
}

// Edges of the Robotstxt.
func (Robotstxt) Edges() []ent.Edge {
	return nil
}
