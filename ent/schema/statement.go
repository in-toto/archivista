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
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Statement represents an in-toto statement from an archived dsse envelope
type Statement struct {
	ent.Schema
}

// Fields of the Statement.
func (Statement) Fields() []ent.Field {
	return []ent.Field{
		field.String("predicate").NotEmpty(),
	}
}

// Edges of the Statement.
func (Statement) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("subjects", Subject.Type).Annotations(entgql.RelayConnection()),
		edge.To("attestation_collections", AttestationCollection.Type).Unique(),

		edge.From("dsse", Dsse.Type).Ref("statement"),
	}
}

func (Statement) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("predicate"),
	}
}
