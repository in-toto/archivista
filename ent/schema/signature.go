package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Signature holds the schema definition for the Signature entity.
type Signature struct {
	ent.Schema
}

// Fields of the Signature.
func (Signature) Fields() []ent.Field {
	return []ent.Field{
		field.String("key_id").NotEmpty(),
		field.String("signature").Unique().NotEmpty(),
	}
}

// Edges of the Signature.
func (Signature) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("dsse", Dsse.Type).Ref("signatures").Unique(),
	}
}
