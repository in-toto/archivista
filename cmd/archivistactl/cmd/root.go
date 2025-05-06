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

package cmd

import (
	"log"
	"net/http"
	"strings"

	"github.com/in-toto/archivista/pkg/api"
	"github.com/spf13/cobra"
)

var (
	archivistaUrl  string
	requestHeaders []string

	rootCmd = &cobra.Command{
		Use:   "archivistactl",
		Short: "A utility to interact with an archivista server",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&archivistaUrl, "archivistaurl", "u", "http://localhost:8082", "url of the archivista instance")
	rootCmd.PersistentFlags().StringArrayVarP(&requestHeaders, "headers", "H", []string{}, "headers to use when making requests to archivista")
}

func Execute() error {
	return rootCmd.Execute()
}

func requestOptions() []api.RequestOption {
	opts := []api.RequestOption{}
	headers := http.Header{}
	for _, header := range requestHeaders {
		headerParts := strings.SplitN(header, ":", 2)
		if len(headerParts) != 2 {
			log.Fatalf("invalid header: %v", header)
		}

		headers.Set(headerParts[0], headerParts[1])
	}

	if len(headers) > 0 {
		opts = append(opts, api.WithHeaders(headers))
	}

	return opts
}
