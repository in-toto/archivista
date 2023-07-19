package cmd

import (
	"context"
	"testing"

	"entgo.io/ent/dialect/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// This test just makes sure our sqllite3 in-memory db works for our tests
func TestJudgeApiServer_WithSQLite3(t *testing.T) {
	// Arrange: Prepare the context and SQL driver
	ctx := context.Background()
	drv, err := sql.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	require.NoError(t, err)

	// Act: Setup the database with the given SQL driver
	judge := SetupDb(ctx, drv)

	// Assert: Check if the server is properly initiated
	require.NotNil(t, judge.authProvider)
	require.NotNil(t, judge.authMiddleware)
	require.NotNil(t, judge.srv)
	require.NotNil(t, judge.database)
	require.NotNil(t, judge.mysqlStoreCh)
	require.NoError(t, err)
}
