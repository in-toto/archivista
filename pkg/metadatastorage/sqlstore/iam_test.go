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
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/stretchr/testify/suite"
)

const (
	mysqlConnStr    = "mysql://user:password@host:3306/dbname"
	postgresConnStr = "postgresql://user:password@host:5432/dbname"
)

type RewriteConnectionStringForIAMSuite struct {
	suite.Suite
}

func TestRewriteConnectionStringForIAM(t *testing.T) {
	suite.Run(t, new(RewriteConnectionStringForIAMSuite))
}

type mockAWSConfigAPI struct {
	cfg aws.Config
	err error
}

func (m *mockAWSConfigAPI) LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	return m.cfg, m.err
}

type mockAWSAuthAPI struct {
	token string
	err   error
}

func (m *mockAWSAuthAPI) BuildAuthToken(ctx context.Context, endpoint, region, dbUser string, creds aws.CredentialsProvider, optFns ...func(options *auth.BuildAuthTokenOptions)) (string, error) {
	return m.token, m.err
}

func (s *RewriteConnectionStringForIAMSuite) TestUnsupportedSqlBackend() {
	AwsConfigAPI = nil
	AwsAuthAPI = nil
	_, err := RewriteConnectionStringForIAM("unsupported", "connstr")
	s.Require().Error(err)
}

func (s *RewriteConnectionStringForIAMSuite) TestMySqlSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		token: "authtoken",
	}
	connStr, err := RewriteConnectionStringForIAM("mysql_rds_iam", mysqlConnStr)
	s.Require().NoError(err)
	s.Equal("mysql://user:authtoken@host:3306/dbname?allowCleartextPasswords=true&tls=true", connStr)
}

func (s *RewriteConnectionStringForIAMSuite) TestMySqlNoPasswordSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		token: "authtoken",
	}
	connStr, err := RewriteConnectionStringForIAM("mysql_rds_iam", "mysql://user@host:3306/dbname")
	s.Require().NoError(err)
	s.Equal("mysql://user:authtoken@host:3306/dbname?allowCleartextPasswords=true&tls=true", connStr)
}

func (s *RewriteConnectionStringForIAMSuite) TestMySqlWithParamsSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		token: "authtoken",
	}
	connStr, err := RewriteConnectionStringForIAM("mysql_rds_iam", fmt.Sprintf("%s?foo=bar", mysqlConnStr))
	s.Require().NoError(err)
	s.Equal("mysql://user:authtoken@host:3306/dbname?allowCleartextPasswords=true&foo=bar&tls=true", connStr)
}

func (s *RewriteConnectionStringForIAMSuite) TestMySqlOverrideParamsSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		token: "authtoken",
	}
	connStr, err := RewriteConnectionStringForIAM("mysql_rds_iam", fmt.Sprintf("%s?allowCleartextPasswords=false&foo=baz&tls=false", mysqlConnStr))
	s.Require().NoError(err)
	s.Equal("mysql://user:authtoken@host:3306/dbname?allowCleartextPasswords=true&foo=baz&tls=true", connStr)
}

func (s *RewriteConnectionStringForIAMSuite) TestPostgresSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		token: "authtoken",
	}
	connStr, err := RewriteConnectionStringForIAM("psql_rds_iam", postgresConnStr)
	s.Require().NoError(err)
	s.Equal("postgresql://user:authtoken@host:5432/dbname", connStr)
}

func (s *RewriteConnectionStringForIAMSuite) TestPostgresNoPasswordSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		token: "authtoken",
	}
	connStr, err := RewriteConnectionStringForIAM("psql_rds_iam", "postgresql://user@host:5432/dbname")
	s.Require().NoError(err)
	s.Equal("postgresql://user:authtoken@host:5432/dbname", connStr)
}

func (s *RewriteConnectionStringForIAMSuite) TestPostgresWithParamsSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		token: "authtoken",
	}
	connStr, err := RewriteConnectionStringForIAM("psql_rds_iam", fmt.Sprintf("%s?foo=bar", postgresConnStr))
	s.Require().NoError(err)
	s.Equal("postgresql://user:authtoken@host:5432/dbname?foo=bar", connStr)
}

func (s *RewriteConnectionStringForIAMSuite) TestLoadConfigError() {
	AwsConfigAPI = &mockAWSConfigAPI{
		err: errors.New("some error"),
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		token: "authtoken",
	}
	_, err := RewriteConnectionStringForIAM("mysql_rds_iam", mysqlConnStr)
	s.Require().Error(err)
	s.Contains(err.Error(), "loading AWS config")
}

func (s *RewriteConnectionStringForIAMSuite) TestBuildAuthTokenError() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		err: errors.New("some error"),
	}
	_, err := RewriteConnectionStringForIAM("mysql_rds_iam", mysqlConnStr)
	s.Require().Error(err)
	s.Contains(err.Error(), "building auth token")
}

func (s *RewriteConnectionStringForIAMSuite) TestMySqlNoHost() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		token: "authtoken",
	}
	_, err := RewriteConnectionStringForIAM("mysql_rds_iam", "mysql://user:password@/dbname")
	s.Require().Error(err)
	s.Contains(err.Error(), "onnection string is missing a host")
}

func (s *RewriteConnectionStringForIAMSuite) TestMySqlNoUser() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		token: "authtoken",
	}
	_, err := RewriteConnectionStringForIAM("mysql_rds_iam", "mysql://:password@host:3306/dbname")
	s.Require().Error(err)
	s.Contains(err.Error(), "onnection string is missing a user")
}

func (s *RewriteConnectionStringForIAMSuite) TestPostgresNoHost() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		token: "authtoken",
	}
	_, err := RewriteConnectionStringForIAM("psql_rds_iam", "postgresql://user:password@/dbname")
	s.Require().Error(err)
	s.Contains(err.Error(), "onnection string is missing a host")
}

func (s *RewriteConnectionStringForIAMSuite) TestPostgresNoUser() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		token: "authtoken",
	}
	_, err := RewriteConnectionStringForIAM("psql_rds_iam", "postgresql://:password@host:5432/dbname")
	s.Require().Error(err)
	s.Contains(err.Error(), "onnection string is missing a user")
}

func (s *RewriteConnectionStringForIAMSuite) TestUrlParseError() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}
	AwsAuthAPI = &mockAWSAuthAPI{
		err: errors.New("some error"),
	}
	_, err := RewriteConnectionStringForIAM("mysql_rds_iam", "://invalid-url")
	s.Require().Error(err)
	s.Contains(err.Error(), "parsing connection string")
}
