package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/mixin"
	"github.com/testifysec/judge/judge-api/ent/privacy"
	"github.com/testifysec/judge/judge-api/rule"
)

type TenantMixin struct {
	mixin.Schema
}

func (TenantMixin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tenant", Tenant.Type).
			Unique().
			Required().
			Immutable(),
	}
}

func (TenantMixin) Policy() ent.Policy {
	return privacy.Policy{
		Mutation: privacy.MutationPolicy{
			rule.DenyIfNoViewer(),
		},
		Query: privacy.QueryPolicy{
			rule.DenyIfNoViewer(),
			rule.DenyIfNoTenants(),
			rule.FilterAccessibleTenants(),
		},
	}
}
