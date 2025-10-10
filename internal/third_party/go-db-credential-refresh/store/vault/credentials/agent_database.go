package vaultcredentials

import (
	"context"
	"os"

	"github.com/davepgreene/go-db-credential-refresh/store"
	"github.com/hashicorp/vault-client-go"
)

// Map handles mapping data from a file on disk to a Credentials object. This
// allows consumers to define how their credential data is structured.
type Mapper func(s string) (*store.Credential, error)

// AgentDatabaseCredentials gets DB credentials the Vault Agent creates on disk
// See: https://www.vaultproject.io/docs/agent/index.html
// One of the key features of the Vault agent is that it can spit out credentials
// using Consul template markup. See https://www.vaultproject.io/docs/agent/template
// for details.
type AgentDatabaseCredentials struct {
	mapper Mapper
	path   string
}

// NewAgentDatabaseCredentials creates a new AgentDatabaseCredentials instance.
func NewAgentDatabaseCredentials(mapper Mapper, path string) CredentialLocation {
	return &AgentDatabaseCredentials{
		mapper: mapper,
		path:   path,
	}
}

// GetCredentials implements the CredentialLocation interface.
func (adb *AgentDatabaseCredentials) GetCredentials(_ context.Context, _ *vault.Client) (string, error) {
	creds, err := os.ReadFile(adb.path)
	if err != nil {
		return "", err
	}

	return string(creds), nil
}

// Map implements the CredentialLocation interface.
func (adb *AgentDatabaseCredentials) Map(s string) (*store.Credential, error) {
	m, err := adb.mapper(s)
	if err != nil {
		return nil, err
	}

	return m, nil
}
