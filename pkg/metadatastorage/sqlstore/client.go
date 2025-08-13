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
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"ariga.io/sqlcomment"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
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
		return nil, fmt.Errorf("unknown sql backend: %s", upperSqlBackend)
	}

	// if upperSqlBackend ends with _RDS_IAM, then rewrite the connection string to use
	// AWS RDS IAM authentication
	if strings.HasSuffix(upperSqlBackend, "_RDS_IAM") {
		var err error
		connectionString, err = rewriteConnectionStringForIAM(sqlBackend, connectionString)
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

// rewriteConnectionStringForIAM rewrites the connection string to use AWS RDS IAM authentication
// It supports both MYSQL_RDS_IAM and PSQL_RDS_IAM backends
func rewriteConnectionStringForIAM(sqlBackend string, connectionString string) (string, error) {
	upperSqlBackend := strings.ToUpper(sqlBackend)
	nURL, err := url.Parse(connectionString)
	if err != nil {
		return "", fmt.Errorf("parsing connection string: %w", err)
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", fmt.Errorf("loading AWS config: %w", err)
	}
	// generate a new rds session auth tokenized connection string
	rdsEndpoint := fmt.Sprintf("%s:%s", nURL.Hostname(), nURL.Port())
	token, err := auth.BuildAuthToken(context.Background(), rdsEndpoint, cfg.Region, nURL.User.Username(), cfg.Credentials)
	if err != nil {
		return "", fmt.Errorf("building auth token: %w", err)
	}
	nURL.User = url.UserPassword(nURL.User.Username(), token)
	// for mysql, we need to add some query parameters
	if strings.HasPrefix(upperSqlBackend, "MYSQL") {
		q := nURL.Query()
		q.Set("tls", "true")
		q.Set("allowCleartextPasswords", "true")
		nURL.RawQuery = q.Encode()
	}
	return nURL.String(), nil
}
