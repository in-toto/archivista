// Copyright 2022 The Archivist Contributors
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
	"github.com/kelseyhightower/envconfig"
	"net/url"
)

type Config struct {
	EnableSPIFFE               bool    `default:"TRUE" desc:"Enable SPIFFE support" split_words:"true"`
	ListenOn                   url.URL `default:"unix:///listen.on.socket" desc:"url to listen on" split_words:"true"`
	LogLevel                   string  `default:"INFO" desc:"Log level" split_words:"true"`
	SPIFFEAddress              string  `default:"unix:///tmp/spire-agent/public/api.sock" desc:"SPIFFE server address" split_words:"true"`
	SPIFFETrustedServerId      string  `default:"" desc:"Trusted SPIFFE server ID; defaults to any" split_words:"true"`
	SQLStoreConnectionString   string  `default:"root:example@tcp(db)/testify" desc:"SQL store connection string" split_words:"true"`
	BlobStoreEndpoint          string  `default:"127.0.0.1:9000" desc:"URL endpoint for blob storage" split_words:"true"`
	BlobStoreAccessKeyId       string  `default:"Blob store access key id" desc:"" split_words:"true"`
	BlobStoreSecretAccessKeyId string  `default:"Blob store secret access key id" desc:"" split_words:"true"`
	BlobStoreUseSSL            bool    `default:"TRUE" desc:"Minio SSL toggle" split_words:"true"`
	BlobStoreBucketName        string  `default:"" desc:"Blob store bucket name" split_words:"true"`
}

// Process reads config from env
func (c *Config) Process() error {
	if err := envconfig.Usage("archivist", c); err != nil {
		return err
	}
	return envconfig.Process("archivist", c)
}
