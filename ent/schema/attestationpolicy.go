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
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Attestation represents an attestation from a witness attestation collection
type AttestationPolicy struct {
	ent.Schema
}

func (AttestationPolicy) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable().Unique(),
		field.String("name").NotEmpty(),
	}
}

// Edges of the AttestationPolicy.
func (AttestationPolicy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("statement", Statement.Type).
			Ref("policy").Unique(),
	}
}

func (AttestationPolicy) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
	}
}

func (AttestationPolicy) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
	}
}
