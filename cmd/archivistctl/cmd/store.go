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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/testifysec/archivist-api/pkg/api/archivist"
	"github.com/testifysec/archivist/internal/server"
	"github.com/testifysec/go-witness/dsse"
)

var (
	storeCmd = &cobra.Command{
		Use:          "store",
		Short:        "stores an attestation on the archivist server",
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := newConn(archivistUrl)
			if err != nil {
				return err
			}

			defer conn.Close()
			for _, filePath := range args {
				if gitoid, err := storeAttestationByPath(cmd.Context(), archivist.NewCollectorClient(conn), filePath); err != nil {
					return fmt.Errorf("failed to store %s: %w", filePath, err)
				} else {
					fmt.Printf("%s stored with gitoid %s\n", filePath, gitoid)
				}
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(storeCmd)
}

func storeAttestationByPath(ctx context.Context, client archivist.CollectorClient, path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()
	return storeAttestation(ctx, client, file)
}

func storeAttestation(ctx context.Context, client archivist.CollectorClient, envelope io.Reader) (string, error) {
	objBytes, err := io.ReadAll(envelope)
	if err != nil {
		return "", err
	}

	obj := &dsse.Envelope{}
	if err := json.Unmarshal(objBytes, &obj); err != nil {
		return "", err
	}

	if len(obj.Payload) == 0 || obj.PayloadType == "" || len(obj.Signatures) == 0 {
		return "", fmt.Errorf("obj is not DSSE %d %d %d", len(obj.Payload), len(obj.PayloadType), len(obj.Signatures))
	}

	stream, err := client.Store(ctx)
	if err != nil {
		return "", err
	}

	chunk := &archivist.Chunk{}
	for curr := 0; curr < len(objBytes); curr += server.ChunkSize {
		if curr+server.ChunkSize > len(objBytes) {
			chunk.Chunk = objBytes[curr:]
		} else {
			chunk.Chunk = objBytes[curr : curr+server.ChunkSize]
		}

		if err := stream.Send(chunk); err != nil {
			return "", err
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	return resp.GetGitoid(), nil
}
