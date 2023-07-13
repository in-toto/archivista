package schema_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testifysec/judge/judge-api/ent"
	"github.com/testifysec/judge/judge-api/ent/privacy"
	"github.com/testifysec/judge/judge-api/ent/user"
	"github.com/testifysec/judge/judge-api/viewer"

	_ "github.com/mattn/go-sqlite3"
)

func TestTenantPermissions(t *testing.T) {
	client := createTestClient(t)
	adminUserID := "root"

	as := bootstrapTestEnvironment(t, client, adminUserID)

	t.Run("Root admin access and permissions", func(t *testing.T) {
		testRootAdminAccessAndPermissions(t, client, adminUserID, as)
	})

	client = createTestClient(t)
	adminUserID = "root"

	as = bootstrapTestEnvironment(t, client, adminUserID)

	t.Run("Non-root admin access and permissions", func(t *testing.T) {
		testNonRootAdminAccessAndPermissions(t, client, adminUserID, as)
	})
}

func testRootAdminAccessAndPermissions(t *testing.T, client *ent.Client, adminUserID string, as AuthStub) {
	v := getViewerFromUser(t, client, adminUserID, as)
	vctx := viewer.NewContext(context.Background(), v)

	// Test case 1: Root admin has access to all tenants
	tenants, err := client.Tenant.Query().All(vctx)
	require.NoError(t, err, "Root admin should be able to query all tenants")
	fmt.Println("Number of tenants: ", len(tenants))

	// Test case 2: Root admin can modify tenant details
	modifiedTenantName := "ModifiedTenant"
	tenantToModify := tenants[0] // Assuming there is at least one tenant

	_, err = client.Tenant.UpdateOneID(tenantToModify.ID).SetName(modifiedTenantName).Save(vctx)
	require.NoError(t, err, "Root admin should be able to modify tenant details")

	// Test case 3: Root admin can delete tenants
	tenantToDelete := tenants[1] // Assuming there are at least two tenants

	v = getViewerFromUser(t, client, "root", as)
	vctx = viewer.NewContext(context.Background(), v)

	err = client.Tenant.DeleteOneID(tenantToDelete.ID).Exec(vctx)
	require.NoError(t, err, "Root admin should be able to delete tenants")
}

func testNonRootAdminAccessAndPermissions(t *testing.T, client *ent.Client, adminUserID string, as AuthStub) {
	v := getViewerFromUser(t, client, "Alice", as)
	alicectx := viewer.NewContext(context.Background(), v)

	// Test case 1: Non-root admin (Alice) can access only assigned tenants
	aliceTenantIDs := v.TenantsAssigned

	_, err := client.Tenant.Query().All(alicectx)
	require.NoError(t, err, "Non-root admin should be able to query assigned tenants")

	_, err = client.Tenant.Get(alicectx, aliceTenantIDs[0]) // Assuming Alice is assigned to at least one tenant
	require.NoError(t, err, "Non-root admin should be able to access assigned tenants")

	// Test case 2: Non-root admin (Alice) cannot access unassigned tenants
	allow := privacy.DecisionContext(context.Background(), privacy.Allow)
	bobID := client.User.Query().Where(user.IdentityID("Bob")).OnlyX(allow).IdentityID
	bobTenantIDs := getUserTenants(t, bobID, as)

	_, err = client.Tenant.Get(alicectx, bobTenantIDs[1]) // Assuming there is at least one unassigned tenant
	require.Error(t, err, "Non-root admin should not be able to access unassigned tenants")

	// Test case 3: Non-root admin (Alice) can modify assigned tenant details but not unassigned tenants
	modifiedAssignedTenantName := "ModifiedAssignedTenant"
	_, err = client.Tenant.UpdateOneID(aliceTenantIDs[0]).SetName(modifiedAssignedTenantName).Save(alicectx)
	require.NoError(t, err, "Non-root admin should be able to modify assigned tenant details")

	// Test case 4: Non-root admin (Alice) cannot modify unassigned tenant details
	modifiedUnassignedTenantName := "ModifiedUnassignedTenant"
	_, err = client.Tenant.UpdateOneID(bobTenantIDs[1]).SetName(modifiedUnassignedTenantName).Save(alicectx)
	require.Error(t, err, "Non-root admin should not be able to modify unassigned tenant details")

	// Test case 5: Non-root admin (Alice) cannot delete tenants they do not own
	err = client.Tenant.DeleteOneID(bobTenantIDs[0]).Exec(alicectx)
	require.Error(t, err, "Non-root admin should not be able to delete tenants")
}
