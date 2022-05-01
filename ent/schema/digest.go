package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Digest holds the schema definition for the Digest entity.
type Digest struct {
	ent.Schema
}

// Fields of the Digest.
func (Digest) Fields() []ent.Field {
	return []ent.Field{
		field.String("algorithm").NotEmpty(),
		field.String("value").NotEmpty(),
	}
}

// Edges of the Digest.
func (Digest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("subject", Subject.Type).
			Ref("digests").
			Unique(),
	}
}
