package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Project struct {
	ent.Schema
}

func (Project) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TenantMixin{},
		MetaDataMixin{},
	}
}

func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.String("repo_id").NotEmpty().Unique(),
		field.String("name").NotEmpty(),
		field.String("projecturl").NotEmpty(),
	}
}

func (Project) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Project) Annotations() []schema.Annotation {
	return []schema.Annotation{
		//entgql.RelayConnection(),
		entgql.QueryField(),
	}
}
