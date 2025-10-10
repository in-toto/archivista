package vaultcredentials

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

func TestGetFromVaultSecretsAPI(t *testing.T) {
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

	// Valid path with response
	b, err := GetFromVaultSecretsAPI(ctx, client.Client, "", "auth/token/lookup-self")
	if err != nil {
		t.Fatal(err)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(b), &resp); err != nil {
		t.Fatal(err)
	}
	// Testing secret data attributes is a bit brittle unfortunately :(
	for _, field := range []string{
		"accessor",
		"creation_time",
		"creation_ttl",
		"id",
		"path",
		"ttl",
		"type",
	} {
		if _, ok := resp[field]; !ok {
			t.Fatalf("expected '%s' to be in response data", field)
		}
	}

	// Invalid path
	_, err = GetFromVaultSecretsAPI(ctx, client.Client, "", "flerp/derp/herp")
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	var respErr *vault.ResponseError
	if !errors.As(err, &respErr) {
		t.Fatalf("expected a '%T' but got '%T' instead", respErr, err)
	}
}

func TestGetFromVaultSecretsAPIWithVaultError(t *testing.T) {
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

	if _, err := client.Client.Secrets.KvV2Write(ctx, "foo", schema.KvV2WriteRequest{
		Data: map[string]interface{}{
			"secret": "string",
		},
	}, vault.WithMountPath("secret")); err != nil {
		t.Fatal(err)
	}

	if _, err := client.Client.Write(ctx, "sys/policy/restricted", map[string]interface{}{
		"policy": `path "secret/foo" {
			capabilities = ["deny"]
		}`,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := client.Client.Auth.TokenCreate(ctx, schema.TokenCreateRequest{
		Policies: []string{"restricted"},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Client.SetToken(resp.Auth.ClientToken); err != nil {
		t.Fatal(err)
	}

	if resp, err := GetFromVaultSecretsAPI(ctx, client.Client, "secret", "foo"); err == nil {
		t.Fatalf("expected an error but got '%s' as a response instead", resp)
	}
}
