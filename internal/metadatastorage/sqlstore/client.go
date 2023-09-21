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
	"strings"
	"time"

	"ariga.io/sqlcomment"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/testifysec/archivista/ent"

	_ "github.com/lib/pq"
)

// NewEntClient creates an ent client for use in the sqlmetadata store.
// Valid backends are MYSQL and PSQL.
func NewEntClient(sqlBackend string, connectionString string) (*ent.Client, error) {
	var entDialect string
	switch strings.ToUpper(sqlBackend) {
	case "MYSQL":
		dbConfig, err := mysql.ParseDSN(connectionString)
		if err != nil {
			return nil, fmt.Errorf("could not parse mysql connection string: %w", err)
		}

		// this tells the go-sql-driver to parse times from mysql to go's time.Time
		// see https://github.com/go-sql-driver/mysql#timetime-support for details
		dbConfig.ParseTime = true
		entDialect = dialect.MySQL
		connectionString = dbConfig.FormatDSN()
	case "PSQL":
		entDialect = dialect.Postgres
	default:
		return nil, fmt.Errorf("unknown sql backend: %s", sqlBackend)
	}

	drv, err := sql.Open(entDialect, connectionString)
	if err != nil {
		return nil, fmt.Errorf("could not open sql connection: %w", err)
	}

	db := drv.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(3 * time.Minute)
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
