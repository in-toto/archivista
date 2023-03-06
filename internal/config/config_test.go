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

package config

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_Process(t *testing.T) {
	// Set up test environment variables
	os.Setenv("ARCHIVISTA_LISTEN_ON", "tcp://0.0.0.0:8082")
	os.Setenv("ARCHIVISTA_LOG_LEVEL", "DEBUG")
	os.Setenv("ARCHIVISTA_CORS_ALLOW_ORIGINS", "http://localhost,https://example.com")
	os.Setenv("ARCHIVISTA_ENABLE_SPIFFE", "FALSE")
	os.Setenv("ARCHIVISTA_SQL_STORE_CONNECTION_STRING", "root:password@tcp(localhost:3306)/testify")
	os.Setenv("ARCHIVISTA_STORAGE_BACKEND", "BLOB")
	os.Setenv("ARCHIVISTA_BLOB_STORE_ENDPOINT", "https://s3.example.com")
	os.Setenv("ARCHIVISTA_BLOB_STORE_ACCESS_KEY_ID", "ACCESSKEYID")
	os.Setenv("ARCHIVISTA_BLOB_STORE_SECRET_ACCESS_KEY_ID", "SECRETACCESSKEYID")
	os.Setenv("ARCHIVISTA_BLOB_STORE_BUCKET_NAME", "mybucket")
	os.Setenv("ARCHIVISTA_GRAPHQL_WEB_CLIENT_ENABLE", "FALSE")

	// Create a Config struct and call Process()
	c := &Config{}
	err := c.Process()
	require.NoError(t, err)

	// Check that the expected values were read from environment variables
	require.Equal(t, "tcp://0.0.0.0:8082", c.ListenOn)
	require.Equal(t, "DEBUG", c.LogLevel)
	require.Equal(t, []string{"http://localhost", "https://example.com"}, c.CORSAllowOrigins)
	require.False(t, c.EnableSPIFFE)
	require.Equal(t, "root:password@tcp(localhost:3306)/testify", c.SQLStoreConnectionString)
	require.Equal(t, "BLOB", c.StorageBackend)
	require.Equal(t, "https://s3.example.com", c.BlobStoreEndpoint)
	require.Equal(t, "ACCESSKEYID", c.BlobStoreAccessKeyId)
	require.Equal(t, "SECRETACCESSKEYID", c.BlobStoreSecretAccessKeyId)
	require.Equal(t, "mybucket", c.BlobStoreBucketName)
	require.False(t, c.GraphqlWebClientEnable)

	// Clean up environment variables
	os.Unsetenv("ARCHIVISTA_LISTEN_ON")
	os.Unsetenv("ARCHIVISTA_LOG_LEVEL")
	os.Unsetenv("ARCHIVISTA_CORS_ALLOW_ORIGINS")
	os.Unsetenv("ARCHIVISTA_ENABLE_SPIFFE")
	os.Unsetenv("ARCHIVISTA_SQL_STORE_CONNECTION_STRING")
	os.Unsetenv("ARCHIVISTA_STORAGE_BACKEND")
	os.Unsetenv("ARCHIVISTA_FILE_SERVE_ON")
	os.Unsetenv("ARCHIVISTA_FILE_DIR")
	os.Unsetenv("ARCHIVISTA_BLOB_STORE_ENDPOINT")
	os.Unsetenv("ARCHIVISTA_BLOB_STORE_ACCESS_KEY_ID")
	os.Unsetenv("ARCHIVISTA_BLOB_STORE_SECRET_ACCESS_KEY_ID")
	os.Unsetenv("ARCHIVISTA_BLOB_STORE_BUCKET_NAME")
	os.Unsetenv("ARCHIVISTA_ENABLE_GRAPHQL")
	os.Unsetenv("ARCHIVISTA_GRAPHQL_WEB_CLIENT_ENABLE")
}

func TestConfig_DeprecationAndPrecedence(t *testing.T) {
	// Set up test environment variables
	os.Setenv("ARCHIVIST_LISTEN_ON", "tcp://0.0.0.0:8082")
	os.Setenv("ARCHIVIST_LOG_LEVEL", "DEBUG")
	os.Setenv("ARCHIVIST_CORS_ALLOW_ORIGINS", "http://localhost:8080")

	// os.Setenv("ARCHIVISTA_LOG_LEVEL", "INFO")
	// os.Setenv("ARCHIVISTA_CORS_ALLOW_ORIGINS", "http://localhost,https://example.com")

	// Set up log output capturing
	var output bytes.Buffer
	log.SetOutput(&output)
	defer func() { log.SetOutput(os.Stderr) }()

	// Create a Config struct and call Process()
	c := &Config{}
	err := c.Process()
	require.NoError(t, err)

	//check that the non deprecated environment variables work

	// Check that the deprecated variables work
	require.Equal(t, "tcp://0.0.0.0:8082", c.ListenOn)
	require.Equal(t, "DEBUG", c.LogLevel)
	require.Equal(t, []string{"http://localhost:8080"}, c.CORSAllowOrigins)

	// Check that the appropriate error is returned if both old and new prefixes are used
	os.Setenv("ARCHIVIST_LISTEN_ON", "tcp://0.0.0.0:8083")
	os.Setenv("ARCHIVISTA_LISTEN_ON", "tcp://0.0.0.0:8084")
	c = &Config{}
	err = c.Process()
	require.Error(t, err)
	require.Equal(t, "both deprecated and new environment variables are being used. Please use only the new environment variables", err.Error())

	// Check that the deprecated environment variables were logged
	require.Contains(t, output.String(), "Using deprecated environment variable ARCHIVIST_LOG_LEVEL")
	require.Contains(t, output.String(), "Using deprecated environment variable ARCHIVIST_CORS_ALLOW_ORIGINS")

	// Clean up environment variables
	os.Unsetenv("ARCHIVIST_LISTEN_ON")
	os.Unsetenv("ARCHIVIST_LOG_LEVEL")
	os.Unsetenv("ARCHIVISTA_LOG_LEVEL")
	os.Unsetenv("ARCHIVIST_CORS_ALLOW_ORIGINS")
	os.Unsetenv("ARCHIVISTA_CORS_ALLOW_ORIGINS")
	os.Unsetenv("ARCHIVISTA_LISTEN_ON")
}
