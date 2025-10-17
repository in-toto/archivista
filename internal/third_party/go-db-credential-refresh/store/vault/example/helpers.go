package main

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/davepgreene/go-db-credential-refresh/driver"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

// setupVault sets up a Vault server running in Docker then enables the plugins/configs we need for
// this example.
func setupVault(
	ctx context.Context,
	dbHost string,
	dbPort int,
) (*vault.Client, func() error, error) {
	fmt.Println("Creating testcontainers Vault instance")

	client, vaultContainer, err := vaulttest.CreateTestVault(ctx)
	if err != nil {
		if vaultContainer != nil {
			if err := vaultContainer.Terminate(ctx); err != nil {
				return nil, nil, err
			}
		}

		return nil, nil, err
	}

	cleanup := func() error {
		return vaultContainer.Terminate(ctx)
	}

	fmt.Println("Mounting the database backend")
	if _, err := client.Client.System.MountsEnableSecretsEngine(ctx, "database",
		schema.MountsEnableSecretsEngineRequest{
			Type: "database",
		}); err != nil {
		return nil, nil, err
	}

	uri := fmt.Sprintf(
		"postgresql://{{username}}:{{password}}@%s:%d/?sslmode=disable",
		dbHost,
		dbPort,
	)

	fmt.Println("Configuring the postgres database and role")
	// TODO: Switch this to using generated methods once the vault client implements the db-specific options
	if _, err := client.Client.Write(ctx,
		fmt.Sprintf("database/config/%s", dbName),
		map[string]any{
			"plugin_name":    "postgresql-database-plugin",
			"allowed_roles":  role,
			"connection_url": uri,
			"username":       username,
			"password":       password,
		}); err != nil {
		return nil, nil, err
	}

	if _, err := client.Client.Secrets.DatabaseWriteRole(ctx, role, schema.DatabaseWriteRoleRequest{
		DbName: dbName,
		CreationStatements: []string{
			`CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'`,
			`GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}"`,
		},
		DefaultTtl: "2s",
		MaxTtl:     "5s",
	}); err != nil {
		return nil, nil, err
	}

	fmt.Println("Vault has been configured")

	return client.Client, cleanup, nil
}

// setupDb sets up a test postgres database for use in the example.
func setupDb(ctx context.Context) (*postgres.PostgresContainer, error) {
	return postgres.Run(ctx,
		"docker.io/postgres:15.2-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
}

// generateDriverConfig is a convenience method for generating a driver config from a testcontainer db
func generateDriverConfig(
	ctx context.Context,
	dbc *postgres.PostgresContainer,
) (*driver.Config, error) {
	connStr, err := dbc.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	connURL, err := url.Parse(connStr)
	if err != nil {
		return nil, err
	}
	hostParts := strings.Split(connURL.Host, ":")
	port, err := strconv.Atoi(hostParts[1])
	if err != nil {
		return nil, err
	}

	opts := make(map[string]string)
	for k, v := range connURL.Query() {
		opts[k] = v[0]
	}

	return &driver.Config{
		Host: host,
		Port: port,
		DB:   dbName,
		Opts: map[string]string{"sslmode": "disable"},
	}, nil
}
