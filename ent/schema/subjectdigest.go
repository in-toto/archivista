package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// SubjectDigest represents the digests of a subject from an in-toto statement
type SubjectDigest struct {
	ent.Schema
}

// Fields of the Digest.
func (SubjectDigest) Fields() []ent.Field {
	return []ent.Field{
		field.String("algorithm").NotEmpty(),
		field.String("value").NotEmpty(),
	}
}

// Edges of the Digest.
func (SubjectDigest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("subject", Subject.Type).
			Ref("subject_digests").
			Unique(),
	}
}
