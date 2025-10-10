package vaultcredentials

import (
	"context"
	"testing"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

func TestNewKvCredentials(t *testing.T) {
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

	path := "test"
	username := "foo"
	password := "bar"
	secretMountPath := vault.WithMountPath("secret")

	if _, err = client.Client.Secrets.KvV2Write(ctx, path, schema.KvV2WriteRequest{
		Data: map[string]any{
			"username": username,
			"password": password,
		}},
		secretMountPath,
	); err != nil {
		t.Fatal(err)
	}

	kvc := NewKvCredentials("secret", path)
	credStr, err := kvc.GetCredentials(ctx, client.Client)
	if err != nil {
		t.Fatal(err)
	}

	creds, err := kvc.Map(credStr)
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
