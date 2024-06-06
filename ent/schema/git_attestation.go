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

// GitAttestation represents an attestation from a witness attestation collection
type GitAttestation struct {
	ent.Schema
}

func (GitAttestation) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable().Unique(),
		field.String("commit_hash"),
		field.String("author"),
		field.String("author_email"),
		field.String("committer_name"),
		field.String("committer_email"),
		field.String("commit_date"),
		field.String("commit_message"),
		field.JSON("status", []string{}),
		field.String("commit_type"),
		field.String("commit_digest"),
		field.String("signature"),
		field.JSON("parent_hashes", []string{}),
		field.String("tree_hash"),
		field.JSON("refs", []string{}),
		field.JSON("remotes", []string{}),
		//field.JSON("tags", []git.Tag{}), // TODO: Add support for gql marshal/unmarshal then revisit this.
	}
}

func (GitAttestation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("attestation", Attestation.Type).Ref("git_attestation").Unique().Required(),
	}
}

func (GitAttestation) Indexes() []ent.Index {
	return []ent.Index{}
}
