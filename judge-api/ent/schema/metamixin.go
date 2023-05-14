package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type MetaDataMixin struct {
	mixin.Schema
}

func (MetaDataMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).StorageKey("oid").Annotations(entgql.OrderField("ID")),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (MetaDataMixin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("created_by", User.Type).
			Unique().
			Immutable(),
		edge.To("modified_by", User.Type).
			Unique(),
	}
}
