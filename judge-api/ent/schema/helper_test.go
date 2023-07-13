package schema_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testifysec/judge/judge-api/ent"
	"github.com/testifysec/judge/judge-api/ent/enttest"
	"github.com/testifysec/judge/judge-api/ent/privacy"
	"github.com/testifysec/judge/judge-api/ent/tenant"
	"github.com/testifysec/judge/judge-api/ent/user"
	"github.com/testifysec/judge/judge-api/viewer"
)

type AuthStub struct {
	userTenantMappings map[string][]uuid.UUID
}

func NewAuthStub() *AuthStub {
	return &AuthStub{
		userTenantMappings: make(map[string][]uuid.UUID),
	}
}

func (a *AuthStub) ValidateAndGetViewer(ctx context.Context, c string) (viewer.UserViewer, error) {
	v := viewer.UserViewer{}
	tenantIDs, ok := a.userTenantMappings[c]
	if !ok {
		return v, fmt.Errorf("Invalid cookie")
	}

	v.IdentityID = c
	v.TenantsAssigned = tenantIDs
	return v, nil
}

func (a *AuthStub) AddTenantToUser(userID string, tenantID uuid.UUID) {
	a.userTenantMappings[userID] = append(a.userTenantMappings[userID], tenantID)
}

func (a *AuthStub) GetUserTenants(userID string) []uuid.UUID {
	return a.userTenantMappings[userID]
}

func createTestUser(t *testing.T, client *ent.Client, identityID string, tenantIDs []uuid.UUID, as AuthStub) *ent.User {
	t.Helper()
	ctx := context.Background()
	now := time.Now()

	u, err := client.User.Create().SetIdentityID(identityID).SetCreatedAt(now).SetUpdatedAt(now).Save(ctx)
	require.NoError(t, err)

	for _, id := range tenantIDs {
		as.AddTenantToUser(identityID, id)
	}
	return u
}

func createTestClient(t *testing.T) *ent.Client {
	return enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1", enttest.WithOptions(ent.Log(t.Log)))
}

func bootstrapTestEnvironment(t *testing.T, client *ent.Client, adminUserID string) AuthStub {
	as := NewAuthStub()
	t.Helper()

	if adminUserID == "" {
		t.Fatalf("admin user id is required, please set adminUserID")
	}

	allow := privacy.DecisionContext(viewer.NewContext(context.Background(), viewer.UserViewer{
		TenantsAssigned: []uuid.UUID{uuid.New()},
		IdentityID:      adminUserID,
	}), privacy.Allow)

	createTestUser(t, client, adminUserID, []uuid.UUID{}, *as)
	root, err := client.Tenant.Create().SetName("Root").SetType(tenant.TypeROOT).SetDescription("Root").Save(allow)
	require.NoError(t, err)
	as.AddTenantToUser(adminUserID, root.ID)

	orgA := createTestTenant(t, client, "Org A", root.ID, tenant.TypeORG, allow)
	orgB := createTestTenant(t, client, "Org B", root.ID, tenant.TypeORG, allow)
	teamA1 := createTestTenant(t, client, "Team A1", orgA.ID, tenant.TypeTEAM, allow)
	teamB1 := createTestTenant(t, client, "Team B1", orgB.ID, tenant.TypeTEAM, allow)
	teamA2 := createTestTenant(t, client, "Team A2", orgA.ID, tenant.TypeTEAM, allow)

	createTestTenant(t, client, "Team A3", orgA.ID, tenant.TypeTEAM, allow)
	createTestTenant(t, client, "Team B2", orgB.ID, tenant.TypeTEAM, allow)
	createTestTenant(t, client, "Team A4", orgA.ID, tenant.TypeTEAM, allow)

	createTestUser(t, client, "Alice", []uuid.UUID{orgA.ID, teamA2.ID}, *as)
	createTestUser(t, client, "Bob", []uuid.UUID{orgB.ID, teamB1.ID}, *as)
	createTestUser(t, client, "Charlie", []uuid.UUID{orgA.ID, teamA1.ID}, *as)
	createTestUser(t, client, "Emma", []uuid.UUID{orgA.ID}, *as)

	tenants, err := client.Tenant.Query().
		WithParent().
		WithChildren().
		All(allow)
	require.NoError(t, err)

	printTenantRelationships(t, tenants)

	return *as
}

func printTenantRelationships(t *testing.T, tenants []*ent.Tenant) {
	t.Log("Tenant relationships:")
	for _, tn := range tenants {
		parent, err := tn.Edges.ParentOrErr()
		if err != nil && tn.Type != tenant.TypeROOT {
			t.Logf("Error retrieving parent for tenant ID %s: %v\n", tn.ID.String(), err)
		}
		children, err := tn.Edges.ChildrenOrErr()
		if err != nil {
			t.Logf("Error retrieving children for tenant ID %s: %v\n", tn.ID.String(), err)
		}
		t.Logf("  - Tenant ID: %s\n", tn.ID.String())
		if parent != nil {
			t.Logf("    - Parent: %s\n", parent.ID.String())
		}
		if children != nil {
			childIDs := make([]string, len(children))
			for i, child := range children {
				childIDs[i] = child.ID.String()
			}
			t.Logf("    - Children: %s\n", strings.Join(childIDs, ", "))
		}
	}
}

func getViewerFromUser(t *testing.T, client *ent.Client, userID string, as AuthStub) *viewer.UserViewer {
	t.Helper()

	ctx := context.Background()
	allow := privacy.DecisionContext(ctx, privacy.Allow)

	_, err := client.User.Query().Where(user.IdentityID(userID)).Only(allow)
	require.NoError(t, err)

	tenantIds := as.GetUserTenants(userID)

	return &viewer.UserViewer{
		TenantsAssigned: tenantIds,
		IdentityID:      userID,
	}
}

func createTestTenant(t *testing.T, client *ent.Client, name string, parentTenantID uuid.UUID, tenantType tenant.Type, allow context.Context) *ent.Tenant {
	t.Helper()

	testTenant, err := client.Tenant.Create().
		SetName(name).
		SetParentID(parentTenantID).
		SetType(tenantType).
		SetDescription(name).
		Save(allow)

	if err != nil {
		t.Fatalf("failed creating tenant: %v", err)
	}
	return testTenant
}

func getUserTenants(t *testing.T, userID string, as AuthStub) []uuid.UUID {
	t.Helper()
	return as.GetUserTenants(userID)
}
