// Copyright 2024 The Archivista Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Metadata struct {
	ent.Schema
}

func (Metadata) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").NotEmpty().Comment("Key value for the metadata item"),
		field.String("value").NotEmpty().Comment("Value for the metadata item"),
	}
}

func (Metadata) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("key", "value").Unique(),
	}
}

func (Metadata) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("envelope", Dsse.Type).
			Ref("metadata"),
	}
}
