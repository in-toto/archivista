// Copyright 2025 The Archivista Contributors
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
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
	"github.com/jkjell/go-db-credential-refresh/driver"
	"github.com/jkjell/go-db-credential-refresh/store/awsrds"
)

func ConfigFromPostgres(connectionString string) (c *driver.Config, user, password string, err error) {
	dc := driver.Config{
		Retries: 3,
	}
	dc.Formatter = driver.PgFormatter

	dbConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, "", "", fmt.Errorf("could not parse postgresql connection string: %w", err)
	}

	dc.Host = dbConfig.Host
	dc.Port = int(dbConfig.Port)
	dc.DB = dbConfig.Database
	dc.Opts = dbConfig.RuntimeParams
	if dc.Opts == nil {
		dc.Opts = make(map[string]string)
	}
	// ParseConfig truncates sslmode param
	if dbConfig.TLSConfig == nil {
		dc.Opts["sslmode"] = "disable"
	}

	if dbConfig.User == "" {
		return nil, "", "", fmt.Errorf("connection string is missing a user")
	}

	return &dc, dbConfig.User, dbConfig.Password, nil
}

func ConfigFromMySQL(connectionString string) (c *driver.Config, user, password string, err error) {
	dc := driver.Config{
		Retries: 3,
	}
	dc.Formatter = driver.MysqlFormatter

	dbConfig, err := mysql.ParseDSN(connectionString)
	if err != nil {
		return nil, "", "", fmt.Errorf("parsing connection string: %w", err)
	}

	addr := strings.Split(dbConfig.Addr, ":")
	dc.Host = addr[0]
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		return nil, "", "", fmt.Errorf("could not parse mysql port: %w", err)
	}
	dc.Port = port
	dc.DB = dbConfig.DBName

	dc.Opts = dbConfig.Params
	// this tells the go-sql-driver to parse times from mysql to go's time.Time
	// see https://github.com/go-sql-driver/mysql#timetime-support for details
	if dc.Opts == nil {
		dc.Opts = make(map[string]string)
	}
	dc.Opts["parseTime"] = "true"

	if dbConfig.User == "" {
		return nil, "", "", fmt.Errorf("connection string is missing a user")
	}

	return &dc, dbConfig.User, dbConfig.Passwd, nil
}

func AWSRDSStoreFromDriverConfig(dc *driver.Config, user string) (driver.Store, error) {
	awsConfig, err := AwsConfigAPI.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("problem loading AWS config: %w", err)
	}

	rdsEndpoint := fmt.Sprintf("%s:%d", dc.Host, dc.Port)
	config := awsrds.Config{
		Credentials: awsConfig.Credentials,
		Endpoint:    rdsEndpoint,
		User:        user,
		Region:      awsConfig.Region,
	}

	return awsrds.NewStore(&config)
}
