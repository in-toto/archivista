package driver

import (
	"fmt"
	"net/url"
	"sort"

	"github.com/go-sql-driver/mysql"
)

// Formatter takes connection string components and assembles them into an implementation-specific conn string/DSN.
type Formatter func(username string, password string, host string, port int, db string, opts map[string]string) string

// MysqlFormatter formats a connection string for the go-sql-driver/mysql lib
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

// PgKVFormatter formats a connection string in the K/V format.
func PgKVFormatter(username, password, host string, port int, db string, opts map[string]string) string {
	s := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s", username, password, host, port, db)

	// Sort opts because Go's non-deterministic behavior around map ordering fucks up string formatting
	// with maps as inputs
	keys := make([]string, 0, len(opts))
	for k := range opts {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		s = fmt.Sprintf("%s %s=%s", s, k, opts[k])
	}

	return s
}

// PgFormatter formats a connection URI for the pq and pgx lib.
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
