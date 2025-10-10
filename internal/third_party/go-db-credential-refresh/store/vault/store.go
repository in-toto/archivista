package vault

import (
	"context"
	"errors"
	"net/http"

	"github.com/davepgreene/go-db-credential-refresh/driver"
	"github.com/hashicorp/vault-client-go"

	vaultauth "github.com/davepgreene/go-db-credential-refresh/store/vault/auth"
	vaultcredentials "github.com/davepgreene/go-db-credential-refresh/store/vault/credentials"
)

// TokenLocation is an interface describing how to get a Vault token
type TokenLocation interface {
	GetToken(ctx context.Context, client *vault.Client) (string, error)
}

// Store is a Store implementation for HashiCorp Vault.
type Store struct {
	client *vault.Client
	cl     vaultcredentials.CredentialLocation
	tl     TokenLocation
	creds  driver.Credentials
}

// Config contains configuration information.
type Config struct {
	Client             *vault.Client
	TokenLocation      TokenLocation
	CredentialLocation vaultcredentials.CredentialLocation
}

var (
	ErrConfigRequired             = errors.New("config is required")
	ErrCredentialLocationRequired = errors.New("credential location is required")
	ErrClientRequired             = errors.New("client is required")
	ErrTokenLocationRequired      = errors.New("token location is required")
)

// NewStore creates a new Vault-backed store.
func NewStore(c *Config) (*Store, error) {
	if c == nil {
		return nil, ErrConfigRequired
	}

	if c.CredentialLocation == nil {
		return nil, ErrCredentialLocationRequired
	}

	client := c.Client
	if client == nil {
		return nil, ErrClientRequired
	}

	ctx := context.Background()

	if c.TokenLocation == nil {
		// If the token location is nil, we should check if the client already
		// has a token by doing a self-lookup
		resp, err := client.Auth.TokenLookUpSelf(ctx)
		if err != nil {
			var rErr *vault.ResponseError
			if errors.As(err, &rErr) {
				if rErr.StatusCode == http.StatusForbidden {
					return nil, ErrTokenLocationRequired
				}
			}

			return nil, err
		}

		token, ok := resp.Data["id"]
		if !ok {
			return nil, ErrTokenLocationRequired
		}

		c.TokenLocation = vaultauth.NewTokenAuth(token.(string))
	}

	token, err := c.TokenLocation.GetToken(ctx, client)
	if err != nil {
		return nil, err
	}

	if err := client.SetToken(token); err != nil {
		return nil, err
	}

	return &Store{
		client: client,
		tl:     c.TokenLocation,
		cl:     c.CredentialLocation,
	}, nil
}

// Get implements the Store interface.
func (v *Store) Get(ctx context.Context) (driver.Credentials, error) {
	if v.creds != nil {
		return v.creds, nil
	}

	return v.Refresh(ctx)
}

// Refresh implements the store interface.
func (v *Store) Refresh(ctx context.Context) (driver.Credentials, error) {
	credStr, err := v.cl.GetCredentials(ctx, v.client)
	if err != nil {
		return nil, err
	}

	creds, err := v.cl.Map(credStr)
	if err != nil {
		return nil, err
	}

	// Cache the credentials
	v.creds = creds

	return creds, nil
}
