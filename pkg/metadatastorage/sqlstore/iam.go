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
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jkjell/go-db-credential-refresh/driver"

	_ "github.com/lib/pq"
)

type AWSConfigAPI interface {
	// LoadDefaultConfig loads the default AWS configuration.
	LoadDefaultConfig(ctx context.Context, opts ...func(*config.LoadOptions) error) (aws.Config, error)
}

type AWSConfig struct{}

func (c *AWSConfig) LoadDefaultConfig(ctx context.Context, opts ...func(*config.LoadOptions) error) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, opts...)
}

var AwsConfigAPI AWSConfigAPI = &AWSConfig{}

// RewriteConnectionStringForIAM rewrites the connection string to use AWS RDS IAM authentication
// It supports both MYSQL_RDS_IAM and PSQL_RDS_IAM backends
func RewriteConnectionStringForIAM(sqlBackend string, connectionString string, dryRun bool) (string, error) {
	var dc *driver.Config
	var user string
	upperSqlBackend := strings.ToUpper(sqlBackend)

	if strings.HasPrefix(upperSqlBackend, "MYSQL") {
		var err error
		dc, user, _, err = ConfigFromMySQL(connectionString)
		if err != nil {
			return "", fmt.Errorf("could not get driver config from mysql connection string: %w", err)
		}
		// for mysql, we need to add some query parameters
		dc.Opts["tls"] = "true"
		dc.Opts["allowCleartextPasswords"] = "true"
	} else if strings.HasPrefix(upperSqlBackend, "PSQL") {
		var err error
		dc, user, _, err = ConfigFromPostgres(connectionString)
		if err != nil {
			return "", fmt.Errorf("could not get driver config from mysql connection string: %w", err)
		}
	} else {
		return "", fmt.Errorf("unsupported sql backend: %s", sqlBackend)
	}

	s, err := AWSRDSStoreFromDriverConfig(dc, user)
	if err != nil {
		return "", fmt.Errorf("could not create credentials refresh store: %w", err)
	}

	creds, err := s.Get(context.Background())
	if err != nil {
		return "", fmt.Errorf("could not get refreshed credentials: %w", err)
	}

	if creds == nil {
		return "", fmt.Errorf("refreshed credentials are nil")
	}

	var password string
	if dryRun {
		password = "authtoken"
	} else {
		password = creds.GetPassword()
	}
	return dc.Formatter(creds.GetUsername(), password, dc.Host, dc.Port, dc.DB, dc.Opts), nil
}
