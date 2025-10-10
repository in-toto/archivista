package vaultcredentials

import (
	"context"
	"encoding/json"

	"github.com/davepgreene/go-db-credential-refresh/store"
	"github.com/hashicorp/vault-client-go"
)

// KvCredentials implements the CredentialLocation interface.
type KvCredentials struct {
	path      string
	mountPath string
}

// NewKvCredentials retrieves credentials from Vault's K/V store.
func NewKvCredentials(mountPath string, path string) CredentialLocation {
	return &KvCredentials{
		path:      path,
		mountPath: mountPath,
	}
}

// GetCredentials implements the CredentialLocation interface.
func (kv *KvCredentials) GetCredentials(ctx context.Context, client *vault.Client) (string, error) {
	resp, err := client.Secrets.KvV2Read(ctx, kv.path, vault.WithMountPath(kv.mountPath))
	if err != nil {
		return "", err
	}
	// Something in Vault's API would have to be horribly broken for the response
	// not to be marshalable but it's worth error checking it as a matter of habit.
	b, err := json.Marshal(resp.Data.Data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Map implements the CredentialLocation interface.
func (*KvCredentials) Map(s string) (*store.Credential, error) {
	return DefaultMapper(s)
}
