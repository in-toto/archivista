package rule

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gitlab.com/testifysec/judge-platform/judge-api/ent"
	"gitlab.com/testifysec/judge-platform/judge-api/ent/predicate"
	"gitlab.com/testifysec/judge-platform/judge-api/ent/privacy"
	"gitlab.com/testifysec/judge-platform/judge-api/ent/tenant"
	"gitlab.com/testifysec/judge-platform/judge-api/viewer"
)

// This needs to be increased if we add more tenant levels
const MAX_RECURSION_DEPTH = 2

// DenyIfNoViewer returns a rule that denies access if the viewer is missing from the context.
func DenyIfNoViewer() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		view := viewer.FromContext(ctx)
		if view == nil {
			return privacy.Denyf("viewer-context is missing")
		}
		return privacy.Skip // Skip to the next privacy rule (equivalent to returning nil).
	})
}

// DenyIfNoTenants returns a rule that denies access if the viewer has no tenants.
func DenyIfNoTenants() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		view := viewer.FromContext(ctx)
		tenants := view.Tenants()
		if len(tenants) == 0 {
			return privacy.Denyf("viewer has no tenants")
		}
		return privacy.Skip // Skip to the next privacy rule (equivalent to returning nil).
	})
}

// FilterAccessibleTenants is a query rule that filters out tenants the user has no access to.
// It checks for direct access to the tenant, access to the parent tenant, access to the child tenants,
// access to the grandparent tenant, and access to the grandchildren tenants.
func tenantIDsToString(tenantIDs []uuid.UUID) string {
	var ids []string
	for _, id := range tenantIDs {
		ids = append(ids, id.String())
	}
	return strings.Join(ids, ", ")
}

func FilterAccessibleTenants() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		view := viewer.FromContext(ctx)
		accessibleTenantIDs := view.Tenants()
		if len(accessibleTenantIDs) == 0 {
			return privacy.Denyf("no accessible tenants in viewer")
		}

		var orPredicates []predicate.Tenant
		orPredicates = append(orPredicates, tenant.IDIn(accessibleTenantIDs...))

		// Loop through the generations to build nested predicates
		for i := 1; i <= MAX_RECURSION_DEPTH; i++ {
			parentPred := tenant.IDIn(accessibleTenantIDs...)
			childPred := tenant.IDIn(accessibleTenantIDs...)

			// Traverse up the generations
			for j := 0; j < i; j++ {
				parentPred = tenant.HasParentWith(parentPred)
			}
			orPredicates = append(orPredicates, parentPred)

			// Traverse down the generations
			for j := 0; j < i; j++ {
				childPred = tenant.HasChildrenWith(childPred)
			}
			orPredicates = append(orPredicates, childPred)
		}

		switch query := q.(type) {
		case *ent.TenantQuery:
			// Filter tenants based on direct access or access to their ancestors/descendants
			query.Where(tenant.Or(orPredicates...))
		// case *ent.EnvironmentQuery:
		// 	fmt.Printf("Filtering environments for viewer with tenant IDs: %s\n", tenantIDsToString(accessibleTenantIDs))
		// 	// Filter environments based on direct access or access to their ancestors/descendants
		// 	query.Where(environment.HasTenantWith(tenant.Or(orPredicates...)))
		default:
			return privacy.Skip
		}

		return privacy.Skip
	})
}

// AllowIfTenantMatchesOrIsChild returns a rule that allows mutations if the viewer has access to the tenant or its parent tenant.
func AllowIfViewerHasAccessToTenantOrAncestor() privacy.MutationRule {
	return privacy.TenantMutationRuleFunc(func(ctx context.Context, m *ent.TenantMutation) error {
		view := viewer.FromContext(ctx)
		vTenants := view.Tenants()
		if len(vTenants) == 0 {
			return privacy.Denyf("viewer has no tenants")
		}

		tenantIDs, err := m.IDs(ctx)
		if err != nil {
			return privacy.Denyf("failed to get tenant IDs")
		}

		accessibleTenantIDs := make(map[uuid.UUID]bool)
		for _, id := range vTenants {
			accessibleTenantIDs[id] = true
		}

		parentIDs, err := m.Client().Tenant.Query().Where(tenant.IDIn(tenantIDs...)).QueryParent().IDs(ctx)
		if err != nil {
			return privacy.Denyf("failed to get parent tenant IDs")
		}

		for _, id := range tenantIDs {
			if accessibleTenantIDs[id] {
				return privacy.Allow
			}
		}

		for _, parentID := range parentIDs {
			if accessibleTenantIDs[parentID] {
				return privacy.Allow
			}
		}

		return privacy.Denyf("viewer does not have access to tenant")
	})
}
