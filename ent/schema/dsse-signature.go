package schema

import (
	"entgo.io/ent"
)

// Statement holds the schema definition for the Statement entity.
type DsseSignature struct {
	ent.Schema
}

// Fields of the Statement.
func (DsseSignature) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the Statement.
func (DsseSignature) Edges() []ent.Edge {
	return []ent.Edge{
		//edge.From("dsse", Dsse.Type).Ref("signatures"),
	}
	return nil
}
