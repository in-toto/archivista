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
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"

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

type AWSAuthAPI interface {
	// BuildAuthToken builds an authentication token for AWS RDS IAM authentication.
	BuildAuthToken(ctx context.Context, endpoint, region, dbUser string, creds aws.CredentialsProvider, optFns ...func(options *auth.BuildAuthTokenOptions)) (string, error)
}

type AWSAuth struct{}

func (a *AWSAuth) BuildAuthToken(ctx context.Context, endpoint, region, dbUser string, creds aws.CredentialsProvider, optFns ...func(options *auth.BuildAuthTokenOptions)) (string, error) {
	return auth.BuildAuthToken(ctx, endpoint, region, dbUser, creds, optFns...)
}

var AwsAuthAPI AWSAuthAPI = &AWSAuth{}

// RewriteConnectionStringForIAM rewrites the connection string to use AWS RDS IAM authentication
// It supports both MYSQL_RDS_IAM and PSQL_RDS_IAM backends
func RewriteConnectionStringForIAM(sqlBackend string, connectionString string) (string, error) {
	if AwsConfigAPI == nil || AwsAuthAPI == nil {
		return "", fmt.Errorf("AWSConfigAPI and AWSAuthAPI must not be nil")
	}
	upperSqlBackend := strings.ToUpper(sqlBackend)
	nURL, err := url.Parse(connectionString)
	if err != nil {
		return "", fmt.Errorf("parsing connection string: %w", err)
	}
	if nURL.Host == "" {
		return "", fmt.Errorf("connection string is missing a host")
	}
	if nURL.User == nil || nURL.User.Username() == "" {
		return "", fmt.Errorf("connection string is missing a user")
	}
	cfg, err := AwsConfigAPI.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", fmt.Errorf("loading AWS config: %w", err)
	}
	// generate a new rds session auth tokenized connection string
	rdsEndpoint := fmt.Sprintf("%s:%s", nURL.Hostname(), nURL.Port())
	token, err := AwsAuthAPI.BuildAuthToken(context.Background(), rdsEndpoint, cfg.Region, nURL.User.Username(), cfg.Credentials)
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
