// Copyright 2022 The Archivist Contributors
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
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Signature represents signatures on a DSSE envelope
type Signature struct {
	ent.Schema
}

// Fields of the Signature.
func (Signature) Fields() []ent.Field {
	return []ent.Field{
		field.String("key_id").NotEmpty(),
		field.String("signature").NotEmpty().SchemaType(map[string]string{dialect.MySQL: "text"}),
	}
}

// Edges of the Signature.
func (Signature) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("dsse", Dsse.Type).Ref("signatures").Unique(),
	}
}

func (Signature) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("key_id"),
	}
}
