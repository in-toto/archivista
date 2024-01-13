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
	"context"
	"fmt"
	"os"

	"github.com/in-toto/archivista/pkg/api"
	"github.com/spf13/cobra"
)

var (
	storeCmd = &cobra.Command{
		Use:          "store",
		Short:        "stores an attestation on the archivista server",
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, filePath := range args {
				if gitoid, err := storeAttestationByPath(cmd.Context(), archivistaUrl, filePath); err != nil {
					return fmt.Errorf("failed to store %s: %w", filePath, err)
				} else {
					rootCmd.Printf("%s stored with gitoid %s\n", filePath, gitoid)
				}
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(storeCmd)
}

func storeAttestationByPath(ctx context.Context, baseUrl, path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()
	resp, err := api.StoreWithReader(ctx, baseUrl, file)
	if err != nil {
		return "", err
	}

	return resp.Gitoid, nil
}
