package vaulttest

import (
	"context"
	"math/rand"
	"time"

	vaultclient "github.com/hashicorp/vault-client-go"
	"github.com/testcontainers/testcontainers-go/modules/vault"
)

const (
	tokenLength = 52
)

type TokenCarryingClient struct {
	Client *vaultclient.Client
	Token  string
}

func CreateTestVault(ctx context.Context) (*TokenCarryingClient, *vault.VaultContainer, error) {
	rootToken := RandString(tokenLength)

	vaultContainer, err := vault.Run(
		ctx,
		"hashicorp/vault:1.15.4",
		vault.WithToken(rootToken),
		vault.WithInitCommand(
			"auth enable kubernetes", // Enable the kubernetes auth method
		),
	)
	if err != nil {
		return nil, nil, err
	}

	hostAddress, err := vaultContainer.HttpHostAddress(ctx)
	if err != nil {
		return nil, vaultContainer, err
	}

	// Create a client that talks to the server, initially authenticating with
	// the root token.
	client, err := vaultclient.New(
		vaultclient.WithAddress(hostAddress),
		vaultclient.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, vaultContainer, err
	}

	if err := client.SetToken(rootToken); err != nil {
		return nil, vaultContainer, err
	}

	return &TokenCarryingClient{
		Client: client,
		Token:  rootToken,
	}, vaultContainer, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))] //nolint:gosec
	}

	return string(b)
}
