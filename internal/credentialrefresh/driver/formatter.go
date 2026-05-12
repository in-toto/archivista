// Vendored from https://github.com/davepgreene/go-db-credential-refresh (MIT).
// See LICENSE in this directory.
//
// Copyright (c) 2022-2024 Dave Greene
// Copyright (c) 2026 The Archivista Contributors

package driver

import (
	"fmt"
	"net/url"

	"github.com/go-sql-driver/mysql"
)

// Formatter takes connection string components and assembles them into an implementation-specific conn string/DSN.
type Formatter func(username string, password string, host string, port int, db string, opts map[string]string) string

// MysqlFormatter formats a connection string for the go-sql-driver/mysql lib.
// NOTE: Currently only supports TCP connections.
func MysqlFormatter(username, password, host string, port int, db string, opts map[string]string) string {
	cfg := mysql.NewConfig()
	cfg.Addr = fmt.Sprintf("%s:%d", host, port)
	cfg.Net = "tcp"
	cfg.User = username
	cfg.Passwd = password
	cfg.DBName = db
	cfg.Params = opts

	return cfg.FormatDSN()
}

// PgFormatter formats a connection URI for the pgx lib.
func PgFormatter(username, password, host string, port int, db string, opts map[string]string) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(username, password),
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   db,
	}

	if len(opts) == 0 {
		return u.String()
	}

	o := url.Values{}

	for k, v := range opts {
		o.Set(k, v)
	}

	u.RawQuery = o.Encode()

	return u.String()
}
