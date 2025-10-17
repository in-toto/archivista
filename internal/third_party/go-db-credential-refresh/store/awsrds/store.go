package awsrds

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/davepgreene/go-db-credential-refresh/driver"
	"github.com/davepgreene/go-db-credential-refresh/store"
)

var (
	errMissingConfig      = errors.New("config is required")
	errMalformedEndpoint  = errors.New("endpoint must be in the form of 'hostname:port'")
	errMissingCredentials = errors.New("credentials cannot be nil")
)

type errMissingConfigItem struct {
	item string
}

func (e errMissingConfigItem) Error() string {
	return fmt.Sprintf("%s is required", e.item)
}

// Store is a Store implementation for AWS RDS.
// https://aws.amazon.com/premiumsupport/knowledge-center/users-connect-rds-iam/
type Store struct {
	*Config
	creds driver.Credentials
}

// Config contains configuration information.
type Config struct {
	Credentials aws.CredentialsProvider
	Endpoint    string // Endpoint takes the form of host:port
	Region      string
	User        string
}

// NewStore creates a new RDS-backed store.
func NewStore(c *Config) (*Store, error) {
	if c == nil {
		return nil, errMissingConfig
	}

	if c.Endpoint == "" {
		return nil, &errMissingConfigItem{item: "endpoint"}
	}

	if c.Region == "" {
		return nil, &errMissingConfigItem{item: "region"}
	}

	if c.User == "" {
		return nil, &errMissingConfigItem{item: "user"}
	}

	if !strings.HasPrefix(c.Endpoint, "http://") && !strings.HasPrefix(c.Endpoint, "https://") {
		c.Endpoint = "http://" + c.Endpoint
	}

	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return nil, err
	}

	if u.Hostname() == "" {
		return nil, errMalformedEndpoint
	}

	if u.Port() == "" {
		return nil, errMalformedEndpoint
	}

	if c.Credentials == nil {
		return nil, errMissingCredentials
	}

	return &Store{
		Config: c,
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
	token, err := auth.BuildAuthToken(ctx, v.Endpoint, v.Region, v.User, v.Credentials)
	if err != nil {
		return nil, err
	}

	creds := &store.Credential{
		Username: v.User,
		Password: token,
	}

	// Cache the credentials
	v.creds = creds

	return creds, nil
}
