// Vendored from https://github.com/davepgreene/go-db-credential-refresh (MIT).
// See LICENSE in the parent credentialrefresh directory.
//
// Copyright (c) 2022-2024 Dave Greene
// Copyright (c) 2026 The Archivista Contributors

package store

// Credential implements the Credentials interface.
type Credential struct {
	Username string
	Password string
}

// GetUsername implements the Credentials interface.
func (c *Credential) GetUsername() string {
	return c.Username
}

// GetPassword implements the Credentials interface.
func (c *Credential) GetPassword() string {
	return c.Password
}
