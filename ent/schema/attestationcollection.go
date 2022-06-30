package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// AttestationCollection represents a witness attestation collection
type AttestationCollection struct {
	ent.Schema
}

func (AttestationCollection) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
	}
}

func (AttestationCollection) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("attestations", Attestation.Type),
		edge.From("statement", Statement.Type).Ref("attestation_collections").Unique().Required(),
	}
}
