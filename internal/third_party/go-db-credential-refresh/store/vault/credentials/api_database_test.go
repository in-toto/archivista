package vaultcredentials

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/vault-client-go/schema"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

func TestNewAPIDatabaseCredentials(t *testing.T) {
	ctx := context.Background()

	client, vaultContainer, err := vaulttest.CreateTestVault(ctx)
	if err != nil {
		if vaultContainer != nil {
			if err := vaultContainer.Terminate(ctx); err != nil {
				t.Fatal(err)
			}
		}
		t.Fatal(err)
	}
	defer func() {
		if err := vaultContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	// Because this CredentialLocation is agnostic to the location of the actual credentials we can
	// fudge this test by using the k/v secret type rather than building a mock vault plugin,
	// mounting it as a db type, and dealing with vault's complicated "separate binary with gRPC
	// communication" process.
	// Instead we mount the k/v secret type at `database`.
	if _, err := client.Client.System.MountsEnableSecretsEngine(ctx, "database", schema.MountsEnableSecretsEngineRequest{
		Type:          "kv",
		PluginVersion: "2",
	}); err != nil {
		t.Fatal(err)
	}

	role := "postgres"

	if _, err := client.Client.Write(ctx, fmt.Sprintf("database/creds/%s", role), map[string]interface{}{
		"username": username,
		"password": password,
	}); err != nil {
		t.Fatal(err)
	}

	adc := NewAPIDatabaseCredentials(role, "")
	credStr, err := adc.GetCredentials(ctx, client.Client)
	if err != nil {
		t.Fatal(err)
	}

	creds, err := adc.Map(credStr)
	if err != nil {
		t.Fatal(err)
	}

	if creds.GetUsername() != username {
		t.Fatalf("expected username to be %s but got %s", username, creds.GetUsername())
	}

	if creds.GetPassword() != password {
		t.Fatalf("expected password to be %s but got %s instead", password, creds.GetPassword())
	}
}
