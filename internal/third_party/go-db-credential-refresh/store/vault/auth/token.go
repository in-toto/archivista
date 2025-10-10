package vaultauth

import (
	"context"
	"errors"

	"github.com/hashicorp/vault-client-go"
)

// TokenAuth is a pass-through authentication mechanism to set vault tokens directly for
// use by the Vault store.
// NOTE: Token renewal should be handled outside of this library.
type TokenAuth struct {
	token string
}

var (
	ErrUnableToLookupToken = errors.New("unable to lookup token information")
)

// NewTokenAuth creates a new Vault token auth location.
func NewTokenAuth(token string) *TokenAuth {
	return &TokenAuth{
		token: token,
	}
}

// GetToken implements the TokenLocation interface.
func (t *TokenAuth) GetToken(ctx context.Context, client *vault.Client) (string, error) {
	if err := client.SetToken(t.token); err != nil {
		return "", err
	}
	// Before we pass the token back we should call an endpoint it will have access to just to be sure
	resp, err := client.Auth.TokenLookUpSelf(ctx)
	if err != nil {
		return "", err
	}
	// We could hit this branch if Vault's token `lookup-self` path is removed but that's pretty unlikely
	// to happen and if it does I'm sure many other things will have broken well before then.
	if resp == nil {
		return "", ErrUnableToLookupToken
	}

	return t.token, nil
}
