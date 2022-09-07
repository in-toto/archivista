package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Signature represents signatures on a DSSE envelope
type Signature struct {
	ent.Schema
}

// Fields of the Signature.
func (Signature) Fields() []ent.Field {
	return []ent.Field{
		field.String("key_id").NotEmpty(),
		field.String("signature").NotEmpty().SchemaType(map[string]string{dialect.MySQL: "text"}),
	}
}

// Edges of the Signature.
func (Signature) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("dsse", Dsse.Type).Ref("signatures").Unique(),
	}
}
