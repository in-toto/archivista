package vaultcredentials

import (
	"encoding/json"
	"errors"

	"github.com/davepgreene/go-db-credential-refresh/store"
)

var (
	errMissingUserName = errors.New("username not set in credential string")
	errMissingPassword = errors.New("password not set in credential string")
)

// DefaultMapper maps the default username/password structure returned from the Vault API.
func DefaultMapper(s string) (*store.Credential, error) {
	var v map[string]any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return nil, err
	}

	username, ok := v["username"].(string)
	if !ok {
		return nil, errMissingUserName
	}

	password, ok := v["password"].(string)
	if !ok {
		return nil, errMissingPassword
	}

	return &store.Credential{
		Username: username,
		Password: password,
	}, nil
}
