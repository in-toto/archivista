package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Statement holds the schema definition for the Statement entity.
type Statement struct {
	ent.Schema
}

// Fields of the Statement.
func (Statement) Fields() []ent.Field {
	return []ent.Field{
		field.String("statement"),
	}
}

// Edges of the Statement.
func (Statement) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("subjects", Subject.Type),
	}
}
