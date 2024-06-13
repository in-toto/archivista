// Copyright 2022 The Archivista Contributors
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
	"github.com/google/uuid"
)

// Attestation represents an attestation from a witness attestation collection
type Mapping struct {
	ent.Schema
}

func (Mapping) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable().Unique(),
		field.String("path").NotEmpty(),
		field.String("type").NotEmpty(),
		field.String("sha1"),
		field.String("sha256"),
		field.String("gitoidSha1").NotEmpty(),
		field.String("gitoidSha256").NotEmpty(),
	}
}

func (Mapping) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posix", Posix.Type),

		edge.From("omnitrail", Omnitrail.Type).Ref("mappings").Unique().Required(),
	}
}

func (Mapping) Indexes() []ent.Index {
	return []ent.Index{}
}
