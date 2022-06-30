package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Attestation represents an attestation from a witness attestation collection
type Attestation struct {
	ent.Schema
}

func (Attestation) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").NotEmpty(),
	}
}

func (Attestation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("attestation_collection", AttestationCollection.Type).Ref("attestations").Unique().Required(),
	}
}
