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

package cmd

import (
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	archivistUrl string

	rootCmd = &cobra.Command{
		Use:   "archivistctl",
		Short: "A utility to interact with an archivist server",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&archivistUrl, "archivisturl", "u", "localhost:8080", "url of the archivist instance")
}

func Execute() error {
	return rootCmd.Execute()
}

func newConn(url string) (*grpc.ClientConn, error) {
	return grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
