package driver

import (
	"context"
)

// Store represents a mechanism for retrieving Credentials.
type Store interface {
	Get(ctx context.Context) (Credentials, error)
	Refresh(ctx context.Context) (Credentials, error)
}

// Credentials represents an abstraction over a username and password.
type Credentials interface {
	GetUsername() string
	GetPassword() string
}
