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
	"net/url"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/suite"
)

const (
	mysqlConnStr    = "user:password@tcp(host:3306)/dbname"
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

func (s *RewriteConnectionStringForIAMSuite) TestUnsupportedSqlBackend() {
	AwsConfigAPI = nil

	_, err := RewriteConnectionStringForIAM("unsupported", "connstr", false)
	s.Require().Error(err)
}

// An example token for PSQL that is returned will look like this:
// user:host:3306?Action=connect&DBUser=user&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=foo%2F20251028%2Fus-east-1%2Frds-db%2Faws4_request&X-Amz-Date=20251028T171846Z&X-Amz-Expires=900&X-Amz-Security-Token=baz&X-Amz-SignedHeaders=host&X-Amz-Signature=e39c63e53919a3d1cdb96211af8d0a83b7afdd637b5d6788ad9014add900489e@tcp(host:3306)/dbname?allowCleartextPasswords=true&parseTime=true&tls=true
//
// The details of the token that we care about are:
// The token starts with `host:3306`
//
// And has these predictable query parameters:
// Action=connect
// DBUser=user
// X-Amz-Credential=foo/20251028/us-east-1/rds-db/aws4_request
// X-Amz-Expires=900
// X-Amz-Security-Token=baz
//
// The other components are variable, based on the time or AWS SDK implementation details
func validateMySqlAuthToken(s *RewriteConnectionStringForIAMSuite, connStr string, opts ...string) {
	parsedConnStr, err := mysql.ParseDSN(connStr)
	s.Require().NoError(err)

	s.Equal("user", parsedConnStr.User)
	s.Equal("tcp", parsedConnStr.Net)
	s.Equal("host:3306", parsedConnStr.Addr)
	s.Equal("dbname", parsedConnStr.DBName)

	// allowCleartextPasswords=true&parseTime=true&tls=true
	s.True(parsedConnStr.AllowCleartextPasswords)
	s.Equal("true", parsedConnStr.TLSConfig)
	s.True(parsedConnStr.ParseTime)

	for _, opt := range opts {
		kv := strings.Split(opt, "=")
		s.Require().Len(kv, 2, "expected opt %s to be in key=value format", opt)
		val, ok := parsedConnStr.Params[kv[0]]
		s.Require().True(ok, "expected param %s to be present", opt)
		s.Equal(kv[1], val, "expected param %s to be true", opt)
	}

	token := parsedConnStr.Passwd
	s.NotEmpty(token)
	s.True(strings.HasPrefix(token, "host:3306"))

	parsedToken, err := url.Parse(token)
	s.Require().NoError(err)

	params := parsedToken.Query()
	s.Equal("connect", params.Get("Action"))
	s.Equal("user", params.Get("DBUser"))
	s.True(strings.HasPrefix(params.Get("X-Amz-Credential"), "foo"))
	s.Contains(params.Get("X-Amz-Credential"), "us-east-1")
	s.Equal("baz", params.Get("X-Amz-Security-Token"))
}

func (s *RewriteConnectionStringForIAMSuite) TestMySqlSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{
			Region:      "us-east-1",
			Credentials: credentials.NewStaticCredentialsProvider("foo", "bar", "baz"),
		},
	}

	connStr, err := RewriteConnectionStringForIAM("mysql_rds_iam", mysqlConnStr, false)
	s.Require().NoError(err)

	validateMySqlAuthToken(s, connStr)
}

func (s *RewriteConnectionStringForIAMSuite) TestMySqlNoPasswordSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{
			Region:      "us-east-1",
			Credentials: credentials.NewStaticCredentialsProvider("foo", "bar", "baz"),
		},
	}

	connStr, err := RewriteConnectionStringForIAM("mysql_rds_iam", "user@tcp(host:3306)/dbname", false)
	s.Require().NoError(err)

	validateMySqlAuthToken(s, connStr)
}

func (s *RewriteConnectionStringForIAMSuite) TestMySqlWithParamsSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{
			Region:      "us-east-1",
			Credentials: credentials.NewStaticCredentialsProvider("foo", "bar", "baz"),
		},
	}

	connStr, err := RewriteConnectionStringForIAM("mysql_rds_iam", fmt.Sprintf("%s?foo=bar", mysqlConnStr), false)
	s.Require().NoError(err)

	validateMySqlAuthToken(s, connStr, "foo=bar")
}

func (s *RewriteConnectionStringForIAMSuite) TestMySqlOverrideParamsSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{
			Region:      "us-east-1",
			Credentials: credentials.NewStaticCredentialsProvider("foo", "bar", "baz"),
		},
	}

	connStr, err := RewriteConnectionStringForIAM("mysql_rds_iam", fmt.Sprintf("%s?allowCleartextPasswords=false&foo=baz&tls=false", mysqlConnStr), false)
	s.Require().NoError(err)

	validateMySqlAuthToken(s, connStr, "foo=baz")
}

// An example token for PSQL that is returned will look like this:
// user:host:5432?Action=connect&DBUser=user&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=foo%2F20251028%2Fus-east-1%2Frds-db%2Faws4_request&X-Amz-Date=20251028T171846Z&X-Amz-Expires=900&X-Amz-Security-Token=baz&X-Amz-SignedHeaders=host&X-Amz-Signature=e39c63e53919a3d1cdb96211af8d0a83b7afdd637b5d6788ad9014add900489e@tcp(host:3306)/dbname?allowCleartextPasswords=true&parseTime=true&tls=true
//
// The details of the token that we care about are:
// The token starts with `host:5432`
//
// And has these predictable query parameters:
// Action=connect
// DBUser=user
// X-Amz-Credential=foo/20251028/us-east-1/rds-db/aws4_request
// X-Amz-Expires=900
// X-Amz-Security-Token=baz
//
// The other components are variable, based on the time or AWS SDK implementation details
func validatePostgresAuthToken(s *RewriteConnectionStringForIAMSuite, connStr string, opts ...string) {
	parsedConnStr, err := url.Parse(connStr)
	s.Require().NoError(err)
	s.Equal("postgres", parsedConnStr.Scheme)
	s.Equal("user", parsedConnStr.User.Username())
	s.Equal("host", parsedConnStr.Hostname())
	s.Equal("5432", parsedConnStr.Port())
	s.Equal("/dbname", parsedConnStr.Path)

	queryParams := parsedConnStr.Query()
	for _, opt := range opts {
		kv := strings.Split(opt, "=")
		s.Require().Len(kv, 2, "expected opt %s to be in key=value format", opt)
		val := queryParams.Get(kv[0])
		s.Equal(kv[1], val, "expected param %s to be true", opt)
	}

	token, set := parsedConnStr.User.Password()
	s.Require().True(set)
	s.NotEmpty(token)
	s.True(strings.HasPrefix(token, "host:5432"))

	parsedToken, err := url.Parse(token)
	s.Require().NoError(err)

	params := parsedToken.Query()
	s.Equal("connect", params.Get("Action"))
	s.Equal("user", params.Get("DBUser"))
	s.True(strings.HasPrefix(params.Get("X-Amz-Credential"), "foo"))
	s.Contains(params.Get("X-Amz-Credential"), "us-east-1")
	s.Equal("baz", params.Get("X-Amz-Security-Token"))
}
func (s *RewriteConnectionStringForIAMSuite) TestPostgresSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{
			Region:      "us-east-1",
			Credentials: credentials.NewStaticCredentialsProvider("foo", "bar", "baz"),
		},
	}

	connStr, err := RewriteConnectionStringForIAM("psql_rds_iam", postgresConnStr, false)
	s.Require().NoError(err)

	validatePostgresAuthToken(s, connStr)
}

func (s *RewriteConnectionStringForIAMSuite) TestPostgresNoPasswordSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{
			Region:      "us-east-1",
			Credentials: credentials.NewStaticCredentialsProvider("foo", "bar", "baz"),
		},
	}

	connStr, err := RewriteConnectionStringForIAM("psql_rds_iam", "postgresql://user@host:5432/dbname", false)
	s.Require().NoError(err)

	validatePostgresAuthToken(s, connStr)
}

func (s *RewriteConnectionStringForIAMSuite) TestPostgresWithParamsSuccess() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{
			Region:      "us-east-1",
			Credentials: credentials.NewStaticCredentialsProvider("foo", "bar", "baz"),
		},
	}

	connStr, err := RewriteConnectionStringForIAM("psql_rds_iam", fmt.Sprintf("%s?foo=bar", postgresConnStr), false)
	s.Require().NoError(err)

	validatePostgresAuthToken(s, connStr, "foo=bar")
	//s.Equal("postgresql://user:authtoken@host:5432/dbname?foo=bar", connStr)
}

func (s *RewriteConnectionStringForIAMSuite) TestLoadConfigError() {
	AwsConfigAPI = &mockAWSConfigAPI{
		err: errors.New("some error"),
	}

	_, err := RewriteConnectionStringForIAM("mysql_rds_iam", mysqlConnStr, false)
	s.Require().Error(err)
	s.Contains(err.Error(), "loading AWS config")
}

func (s *RewriteConnectionStringForIAMSuite) TestBuildAuthTokenError() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}

	_, err := RewriteConnectionStringForIAM("mysql_rds_iam", mysqlConnStr, false)
	s.Require().Error(err)
	s.Contains(err.Error(), "could not create credentials refresh store")
}

func (s *RewriteConnectionStringForIAMSuite) TestMySqlNoUser() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}

	_, err := RewriteConnectionStringForIAM("mysql_rds_iam", ":password@tcp(host:3306)/dbname", false)
	s.Require().Error(err)
	s.Contains(err.Error(), "onnection string is missing a user")
}

func (s *RewriteConnectionStringForIAMSuite) TestPostgresNoUser() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}

	_, err := RewriteConnectionStringForIAM("psql_rds_iam", "postgresql://:password@host:5432/dbname", false)
	s.Require().Error(err)
	s.Contains(err.Error(), "onnection string is missing a user")
}

func (s *RewriteConnectionStringForIAMSuite) TestUrlParseError() {
	AwsConfigAPI = &mockAWSConfigAPI{
		cfg: aws.Config{Region: "us-east-1"},
	}

	_, err := RewriteConnectionStringForIAM("mysql_rds_iam", "://invalid-url", false)
	s.Require().Error(err)
	s.Contains(err.Error(), "parsing connection string")
}
