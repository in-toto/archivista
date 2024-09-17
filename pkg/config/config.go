// Copyright 2022 The Archivista Contributors
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
	"errors"
	"os"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ListenOn         string   `default:"tcp://127.0.0.1:8082" desc:"URL endpoint for Archivista to listen on" split_words:"true"`
	ReadTimeout      int      `default:"120" desc:"HTTP read timeout in seconds" split_words:"true"`
	WriteTimeout     int      `default:"120" desc:"HTTP write timeout in seconds" split_words:"true"`
	LogLevel         string   `default:"INFO" desc:"Log level" split_words:"true"`
	CORSAllowOrigins []string `default:"" desc:"Comma separated list of origins to allow CORS requests from" split_words:"true"`

	EnableSPIFFE          bool   `default:"TRUE" desc:"*** Enable SPIFFE support" split_words:"true"`
	SPIFFEAddress         string `default:"unix:///tmp/spire-agent/public/api.sock" desc:"SPIFFE server address" split_words:"true"`
	SPIFFETrustedServerId string `default:"" desc:"Trusted SPIFFE server ID; defaults to any" split_words:"true"`

	SQLStoreConnectionString      string        `default:"root:example@tcp(db)/testify" desc:"SQL store connection string" split_words:"true"`
	SQLStoreBackend               string        `default:"MYSQL" desc:"SQL backend to use. Options are MYSQL, PSQL" split_words:"true"`
	SQLStoreMaxIdleConnections    int           `default:"10" desc:"Maximum number of connections in the idle connection pool" split_words:"true"`
	SQLStoreMaxOpenConnections    int           `default:"100" desc:"Maximum number of open connections to the database" split_words:"true"`
	SQLStoreConnectionMaxLifetime time.Duration `default:"3m" desc:"Maximum amount of time a connection may be reused" split_words:"true"`

	StorageBackend             string `default:"" desc:"Backend to use for attestation storage. Options are FILE, BLOB, or empty string for disabled." split_words:"true"`
	FileServeOn                string `default:"" desc:"What address to serve files on. Only valid when using FILE storage backend." split_words:"true"`
	FileDir                    string `default:"/tmp/archivista/" desc:"Directory to store and serve files. Only valid when using FILE storage backend." split_words:"true"`
	BlobStoreEndpoint          string `default:"127.0.0.1:9000" desc:"URL endpoint for blob storage. Only valid when using BLOB storage backend." split_words:"true"`
	BlobStoreCredentialType    string `default:"ACCESS_KEY" desc:"Blob store credential type. Options are IAM or ACCESS_KEY" split_words:"true"`
	BlobStoreAccessKeyId       string `default:"" desc:"Blob store access key id. Only valid when using BLOB storage backend." split_words:"true"`
	BlobStoreSecretAccessKeyId string `default:"" desc:"Blob store secret access key id. Only valid when using BLOB storage backend." split_words:"true"`
	BlobStoreUseTLS            bool   `default:"TRUE" desc:"Use TLS for BLOB storage backend. Only valid when using BLOB storage backend." split_words:"true"`
	BlobStoreBucketName        string `default:"" desc:"Bucket to use for storage.  Only valid when using BLOB storage backend." split_words:"true"`

	EnableGraphql          bool `default:"TRUE" desc:"*** Enable GraphQL Endpoint" split_words:"true"`
	GraphqlWebClientEnable bool `default:"TRUE" desc:"Enable GraphiQL, the GraphQL web client" split_words:"true"`

	EnableArtifactStore bool   `default:"FALSE" desc:"*** Enable Artifact Store Endpoints" split_words:"true"`
	ArtifactStoreConfig string `default:"/tmp/artifacts/config.yaml" desc:"Location of the config describing available artifacts" split_words:"true"`

	Publisher                  []string `default:"" desc:"Publisher to use. Options are DAPR, RSTUF or empty string for disabled." split_words:"true"`
	PublisherDaprHost          string   `default:"http://127.0.0.1" desc:"Host for Dapr" split_words:"true"`
	PublisherDaprPort          string   `default:"3500" desc:"Port for Dapr" split_words:"true"`
	PublisherDaprURL           string   `default:"" desc:"URL for Dapr" split_words:"true"`
	PublisherDaprComponentName string   `default:"archivista" desc:"Dapr pubsub component name" split_words:"true"`
	PublisherDaprTopic         string   `default:"attestations" desc:"Dapr pubsub topic" split_words:"true"`
	PublisherRstufHost         string   `default:"http://127.0.0.1" desc:"Host for RSTUF" split_words:"true"`
}

// Process reads config from env
func (c *Config) Process() error {
	if err := envconfig.Usage("archivista", c); err != nil {
		return err
	}

	// Check if old env variables starting with ARCHIVIST_ are being used
	// TODO: Remove this in a future release
	usingDeprecatedEnv := false
	usingNewEnv := false
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "ARCHIVIST_") {
			usingDeprecatedEnv = true
			logrus.Warnf("Using deprecated environment variable %s. Please use ARCHIVISTA_ instead.", e)
		}
		if strings.HasPrefix(e, "ARCHIVISTA_") {
			usingNewEnv = true
		}
	}

	//check if both are being used and error if so
	if usingDeprecatedEnv && usingNewEnv {
		err := errors.New("both deprecated and new environment variables are being used. Please use only the new environment variables")
		return err
	}

	if usingDeprecatedEnv {
		err := envconfig.Process("archivist", c)
		if err != nil {
			return err
		}
	}

	if usingNewEnv {
		err := envconfig.Process("archivista", c)
		if err != nil {
			return err
		}
	}

	return nil
}
