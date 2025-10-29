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

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/in-toto/archivista/pkg/api"
	"github.com/in-toto/archivista/pkg/sigstorebundle"
	"github.com/spf13/cobra"
)

var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Manage Sigstore bundles",
}

var importBundleCmd = &cobra.Command{
	Use:          "import [file]",
	Short:        "Import a Sigstore bundle",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		if gitoid, err := importBundleByPath(cmd.Context(), archivistaUrl, filePath); err != nil {
			return fmt.Errorf("failed to import bundle %s: %w", filePath, err)
		} else {
			rootCmd.Printf("%s imported with gitoid %s\n", filePath, gitoid)
		}
		return nil
	},
}

var exportBundleCmd = &cobra.Command{
	Use:          "export [dsse-id]",
	Short:        "Export a Sigstore bundle",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dsseID := args[0]
		output, _ := cmd.Flags().GetString("output")

		bundleData, err := exportBundleByID(cmd.Context(), archivistaUrl, dsseID)
		if err != nil {
			return fmt.Errorf("failed to export bundle: %w", err)
		}

		if output == "-" {
			fmt.Print(string(bundleData))
		} else {
			if err := os.WriteFile(output, bundleData, 0o600); err != nil {
				return fmt.Errorf("failed to write bundle to %s: %w", output, err)
			}
			rootCmd.Printf("Bundle exported to %s\n", output)
		}

		return nil
	},
}

func init() {
	bundleCmd.AddCommand(importBundleCmd)
	bundleCmd.AddCommand(exportBundleCmd)

	exportBundleCmd.Flags().StringP("output", "o", "-", "Output file (- for stdout)")

	rootCmd.AddCommand(bundleCmd)
}

func importBundleByPath(ctx context.Context, baseUrl, path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	resp, err := api.StoreWithReader(ctx, baseUrl, file, requestOptions()...)
	if err != nil {
		return "", err
	}

	return resp.Gitoid, nil
}

func exportBundleByID(ctx context.Context, baseUrl, dsseID string) ([]byte, error) {
	// Download the DSSE envelope from the API
	envelope, err := api.Download(ctx, baseUrl, dsseID, requestOptions()...)
	if err != nil {
		return nil, fmt.Errorf("failed to download DSSE envelope: %w", err)
	}

	// Reconstruct a Sigstore bundle from the DSSE envelope
	bundle, err := sigstorebundle.ReconstructBundleFromEnvelope(&envelope)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct bundle: %w", err)
	}

	bundleJSON, err := json.Marshal(bundle)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bundle: %w", err)
	}

	return bundleJSON, nil
}
