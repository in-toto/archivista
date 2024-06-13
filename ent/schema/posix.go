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
type Posix struct {
	ent.Schema
}

func (Posix) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable().Unique(),
		field.String("atime"),
		field.String("ctime"),
		field.String("creation_time"),
		field.String("extended_attributes"),
		field.String("file_device_id"),
		field.String("file_flags"),
		field.String("file_inode"),
		field.String("file_system_id"),
		field.String("file_type"),
		field.String("hard_link_count"),
		field.String("mtime"),
		field.String("metadata_ctime"),
		field.String("owner_gid"),
		field.String("owner_uid"),
		field.String("permissions"),
		field.String("size"),
	}
}

func (Posix) Edges() []ent.Edge {
	return []ent.Edge{
		// edge.To("mappings", Mapping.Type),

		edge.From("mapping", Mapping.Type).Ref("posix").Unique().Required(),
	}
}

func (Posix) Indexes() []ent.Index {
	return []ent.Index{
		// index.Fields("type"),
	}
}
