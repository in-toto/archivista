// Vendored from https://github.com/davepgreene/go-db-credential-refresh (MIT).
// See LICENSE in this directory.
//
// Copyright (c) 2022-2024 Dave Greene
// Copyright (c) 2026 The Archivista Contributors

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
