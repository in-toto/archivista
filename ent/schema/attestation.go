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
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Attestation represents an attestation from a witness attestation collection
type Attestation struct {
	ent.Schema
}

func (Attestation) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").NotEmpty(),
	}
}

func (Attestation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("attestation_collection", AttestationCollection.Type).Ref("attestations").Unique().Required(),
	}
}

func (Attestation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
	}
}
