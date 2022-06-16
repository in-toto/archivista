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
	"net/url"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	EnableSPIFFE bool    `default:"TRUE" desc:"Enable SPIFFE support" split_words:"true"`
	ListenOn     url.URL `default:"unix:///listen.on.socket" desc:"url to listen on" split_words:"true"`
	LogLevel     string  `default:"INFO" desc:"Log level" split_words:"true"`

	FileServeOn string `default:"" desc:"What address to serve files on, leave empty to shut off" split_words:"true"`
	FileDir     string `default:"/tmp/archivist/" desc:"Directory to store and serve files" split_words:"true"`
}

// Process reads config from env
func (c *Config) Process() error {
	if err := envconfig.Usage("archivist", c); err != nil {
		return err
	}
	return envconfig.Process("archivist", c)
}
