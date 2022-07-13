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
			defer conn.Close()
			if err != nil {
				return err
			}

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
	defer file.Close()
	if err != nil {
		return "", err
	}

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
