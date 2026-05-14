// Vendored from https://github.com/davepgreene/go-db-credential-refresh (MIT).
// See LICENSE in the parent credentialrefresh directory.
//
// Copyright (c) 2022-2024 Dave Greene
// Copyright (c) 2026 The Archivista Contributors

package static

import (
	"context"

	"github.com/in-toto/archivista/internal/credentialrefresh/driver"
	"github.com/in-toto/archivista/internal/credentialrefresh/store"
)

type Static struct {
	credentials *store.Credential
}

func NewStaticStore(username, password string) *Static {
	return &Static{
		credentials: &store.Credential{
			Username: username,
			Password: password,
		},
	}
}

func (s Static) Get(_ context.Context) (driver.Credentials, error) {
	return s.credentials, nil
}

func (s Static) Refresh(_ context.Context) (driver.Credentials, error) {
	return s.credentials, nil
}
