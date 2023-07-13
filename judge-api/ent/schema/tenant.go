package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/testifysec/judge/judge-api/ent/privacy"
	"github.com/testifysec/judge/judge-api/rule"
)

type TenantType string

const (
	RootTenant TenantType = "ROOT"
	OrgTenant  TenantType = "ORG"
	TeamTenant TenantType = "TEAM"
)

type Tenant struct {
	ent.Schema
}

func (Tenant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		MetaDataMixin{},
	}
}

func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
		field.Enum("type").Values(string(RootTenant), string(OrgTenant), string(TeamTenant)),
	}
}

func (Tenant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
		edge.To("children", Tenant.Type).
			From("parent").Unique(),
	}
}

func (Tenant) Policy() ent.Policy {
	return privacy.Policy{
		Mutation: privacy.MutationPolicy{
			rule.DenyIfNoViewer(),
			rule.DenyIfNoTenants(),
			rule.AllowIfViewerHasAccessToTenantOrAncestor(),
		},
		Query: privacy.QueryPolicy{
			rule.DenyIfNoViewer(),
			rule.DenyIfNoTenants(),
			rule.FilterAccessibleTenants(),
		},
	}
}
