package vaultauth

import (
	"context"
	"testing"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

func TestTokenAuth(t *testing.T) {
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

	ta := NewTokenAuth(client.Token)
	token, err := ta.GetToken(ctx, client.Client)
	if err != nil {
		t.Fatal(err)
	}

	if token != client.Token {
		t.Fatalf("expected token to be %s but got %s instead", client.Token, token)
	}
}

func TestTokenAuthWithInvalidToken(t *testing.T) {
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

	ta := NewTokenAuth("foobar")
	if _, err := ta.GetToken(ctx, client.Client); err == nil {
		t.Fatal("expected an error but didn't get one")
	}
}
