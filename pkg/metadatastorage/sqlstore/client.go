// Copyright 2023 The Archivista Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlstore

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"ariga.io/sqlcomment"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/in-toto/archivista/ent"

	_ "github.com/lib/pq"
)

type ClientOption func(*clientOptions)

type clientOptions struct {
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
}

// Configures a client with the specified max idle connections. Default is 10 connections
func ClientWithMaxIdleConns(maxIdleConns int) ClientOption {
	return func(co *clientOptions) {
		co.maxIdleConns = maxIdleConns
	}
}

// Configures a client with the specified max open connections. Default is 100 connections
func ClientWithMaxOpenConns(maxOpenConns int) ClientOption {
	return func(co *clientOptions) {
		co.maxOpenConns = maxOpenConns
	}
}

// Congiures a client with the specified max connection lifetime. Default is 3 minutes
func ClientWithConnMaxLifetime(connMaxLifetime time.Duration) ClientOption {
	return func(co *clientOptions) {
		co.connMaxLifetime = connMaxLifetime
	}
}

// ensureMySQLConnectionString ensures the connection string has the tcp protocol as required by the go-sql-driver
func ensureMySQLConnectionString(connStr string) (string, error) {
	schema := "mysql://"

	if strings.Contains(connStr, "@tcp(") {
		return connStr, nil
	}

	// Add mysql:// prefix if not present. URL Parse will fail silently if a schema is not present
	if !strings.HasPrefix(connStr, schema) {
		connStr = schema + connStr
	}

	// Parse the connection string as a URL
	u, err := url.Parse(connStr)
	if err != nil {
		return "", fmt.Errorf("invalid mysql connection string: %w", err)
	}

	// Modify the host to include tcp
	u.Host = "tcp(" + u.Host + ")"

	// Remove the mysql:// prefix from the final string
	result := strings.TrimPrefix(u.String(), schema)
	return result, nil
}

// NewEntClient creates an ent client for use in the sqlmetadata store.
// Valid backends are MYSQL and PSQL.
func NewEntClient(sqlBackend string, connectionString string, opts ...ClientOption) (*ent.Client, error) {
	clientOpts := &clientOptions{
		maxIdleConns:    10,
		maxOpenConns:    100,
		connMaxLifetime: 3 * time.Minute,
	}

	for _, opt := range opts {
		opt(clientOpts)
	}

	var entDialect string
	upperSqlBackend := strings.ToUpper(sqlBackend)
	if strings.HasPrefix(upperSqlBackend, "MYSQL") {
		// Ensure the connection string has the tcp protocol as required by the go-sql-driver
		var err error
		connectionString, err = ensureMySQLConnectionString(connectionString)
		if err != nil {
			return nil, fmt.Errorf("could not ensure mysql connection string: %w", err)
		}
		dbConfig, err := mysql.ParseDSN(connectionString)
		if err != nil {
			return nil, fmt.Errorf("could not parse mysql connection string: %w", err)
		}

		// this tells the go-sql-driver to parse times from mysql to go's time.Time
		// see https://github.com/go-sql-driver/mysql#timetime-support for details
		dbConfig.ParseTime = true
		entDialect = dialect.MySQL
		connectionString = dbConfig.FormatDSN()
	} else if strings.HasPrefix(upperSqlBackend, "PSQL") {
		entDialect = dialect.Postgres
	} else {
		return nil, fmt.Errorf("unknown sql backend: %s", sqlBackend)
	}

	// if upperSqlBackend ends with _RDS_IAM, then rewrite the connection string to use
	// AWS RDS IAM authentication
	if strings.HasSuffix(upperSqlBackend, "_RDS_IAM") {
		var err error
		connectionString, err = RewriteConnectionStringForIAM(sqlBackend, connectionString)
		if err != nil {
			return nil, fmt.Errorf("could not rewrite connection string for IAM: %w", err)
		}
	}

	drv, err := sql.Open(entDialect, connectionString)
	if err != nil {
		return nil, fmt.Errorf("could not open sql connection: %w", err)
	}

	db := drv.DB()
	db.SetMaxIdleConns(clientOpts.maxIdleConns)
	db.SetMaxOpenConns(clientOpts.maxOpenConns)
	db.SetConnMaxLifetime(clientOpts.connMaxLifetime)
	sqlcommentDrv := sqlcomment.NewDriver(drv,
		sqlcomment.WithDriverVerTag(),
		sqlcomment.WithTags(sqlcomment.Tags{
			sqlcomment.KeyApplication: "archivista",
			sqlcomment.KeyFramework:   "net/http",
		}),
	)

	client := ent.NewClient(ent.Driver(sqlcommentDrv))
	return client, nil
}
