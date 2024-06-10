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
	"github.com/google/uuid"
)

// SarifRule represents a Sarif Rule
type SarifRule struct {
	ent.Schema
}

func (SarifRule) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable().Unique(),
		field.String("rule_id").NotEmpty(),
		field.String("rule_name").NotEmpty(),
		field.String("short_description").NotEmpty(),
	}
}

func (SarifRule) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("sarif", Sarif.Type).Ref("sarif_rules").Unique(),
	}
}

func (SarifRule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("rule_id"),
		index.Fields("rule_name"),
		index.Fields("short_description"),
	}
}
