// Vendored from https://github.com/davepgreene/go-db-credential-refresh (MIT).
// See LICENSE in this directory.
//
// Copyright (c) 2022-2024 Dave Greene
// Copyright (c) 2026 The Archivista Contributors

package driver

import (
	"database/sql/driver"
	"errors"
	"strings"
)

// AuthError is a func to evaluate the DB-specific error string that indicates an authentication error.
type AuthError func(e error) bool

const (
	MysqlErrorText = "access denied for user"
	PgErrorText    = "authentication failed for user"
)

// MySQLAuthError tests whether an error from MySQL is an authentication failure.
var MySQLAuthError = errorTester(MysqlErrorText)

// PostgreSQLAuthError tests whether an error from PostgreSQL is an authentication failure.
var PostgreSQLAuthError = errorTester(PgErrorText)

func errorTester(text string) AuthError {
	return func(e error) bool {
		return strings.Contains(strings.ToLower(e.Error()), text) || errors.Is(e, driver.ErrBadConn)
	}
}
