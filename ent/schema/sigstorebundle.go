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
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// SigstoreBundle represents a Sigstore bundle v0.3+ wrapping a DSSE envelope
type SigstoreBundle struct {
	ent.Schema
}

// Fields of the SigstoreBundle.
func (SigstoreBundle) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable().Unique(),
		field.String("gitoid_sha256").Unique(),
		field.String("media_type"),
		field.String("version").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the SigstoreBundle.
func (SigstoreBundle) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("dsse", Dsse.Type).Ref("bundle").Unique(),
	}
}

// Indexes of the SigstoreBundle.
func (SigstoreBundle) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("gitoid_sha256"),
	}
}
