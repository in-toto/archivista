package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Dsse represents some metadata about an archived DSSE envelope
type Dsse struct {
	ent.Schema
}

// Fields of the Statement.
func (Dsse) Fields() []ent.Field {
	return []ent.Field{
		field.String("gitbom_sha256").NotEmpty().Unique(),
		field.String("payload_type").NotEmpty(),
	}
}

// Edges of the Statement.
func (Dsse) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("statement", Statement.Type).Unique(),
		edge.To("signatures", Signature.Type),
		edge.To("payload_digests", PayloadDigest.Type),
	}
}
