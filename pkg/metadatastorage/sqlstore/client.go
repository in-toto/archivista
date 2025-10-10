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
	"database/sql"
	"fmt"
	"strings"
	"time"

	"ariga.io/sqlcomment"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/in-toto/archivista/ent"
	"github.com/jkjell/go-db-credential-refresh/driver"
	"github.com/jkjell/go-db-credential-refresh/store/static"
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

	var s driver.Store
	var dc *driver.Config
	var entDialect, driverName, user, password string

	upperSqlBackend := strings.ToUpper(sqlBackend)
	if strings.HasPrefix(upperSqlBackend, "MYSQL") {
		var err error
		dc, user, password, err = ConfigFromMySQL(connectionString)
		if err != nil {
			return nil, fmt.Errorf("could not get driver config from mysql connection string: %w", err)
		}

		entDialect = dialect.MySQL
		driverName = "mysql"
	} else if strings.HasPrefix(upperSqlBackend, "PSQL") {
		var err error
		dc, user, password, err = ConfigFromPostgres(connectionString)
		if err != nil {
			return nil, fmt.Errorf("could not get driver config from postgres connection string: %w", err)
		}

		entDialect = dialect.Postgres
		driverName = "pgx"
	} else {
		return nil, fmt.Errorf("unknown sql backend: %s", sqlBackend)
	}

	// if upperSqlBackend ends with _RDS_IAM, then use awsrds store
	if strings.HasSuffix(upperSqlBackend, "_RDS_IAM") {
		var err error
		s, err = AWSRDSStoreFromDriverConfig(dc, user)
		if err != nil {
			return nil, fmt.Errorf("could not create credentials refresh store: %w", err)
		}
	} else {
		s = static.NewStaticStore(user, password)
		dc.Retries = 0 // no retries for static credentials
	}

	c, err := driver.NewConnector(s, driverName, dc)
	if err != nil {
		return nil, fmt.Errorf("could not create connector: %w", err)
	}

	db := sql.OpenDB(c)
	db.SetMaxIdleConns(clientOpts.maxIdleConns)
	db.SetMaxOpenConns(clientOpts.maxOpenConns)
	db.SetConnMaxLifetime(clientOpts.connMaxLifetime)

	drv := entsql.OpenDB(entDialect, db)
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
