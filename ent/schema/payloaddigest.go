package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// PayloadDigest represents the digest of the payload of a DSSE envelope
type PayloadDigest struct {
	ent.Schema
}

// Fields of the Digest.
func (PayloadDigest) Fields() []ent.Field {
	return []ent.Field{
		field.String("algorithm").NotEmpty(),
		field.String("value").NotEmpty(),
	}
}

// Edges of the Digest.
func (PayloadDigest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("dsse", Dsse.Type).
			Ref("payload_digests").
			Unique(),
	}
}
