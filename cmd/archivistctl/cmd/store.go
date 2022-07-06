package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/testifysec/archivist-api/pkg/api/archivist"
	"github.com/testifysec/go-witness/dsse"
)

var (
	storeCmd = &cobra.Command{
		Use:          "store",
		Short:        "stores an attestation on the archivist server",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := newConn(archivistUrl)
			defer conn.Close()
			if err != nil {
				return err
			}

			file, err := os.Open(args[0])
			defer file.Close()
			if err != nil {
				return err
			}

			return storeAttestation(cmd.Context(), archivist.NewCollectorClient(conn), file)
		},
	}
)

func init() {
	rootCmd.AddCommand(storeCmd)
}

func storeAttestation(ctx context.Context, client archivist.CollectorClient, envelope io.Reader) error {
	objBytes, err := io.ReadAll(envelope)
	if err != nil {
		return err
	}

	obj := &dsse.Envelope{}
	if err := json.Unmarshal(objBytes, &obj); err != nil {
		return err
	}

	if len(obj.Payload) == 0 || obj.PayloadType == "" || len(obj.Signatures) == 0 {
		return fmt.Errorf("obj is not DSSE %d %d %d", len(obj.Payload), len(obj.PayloadType), len(obj.Signatures))
	}

	if _, err := client.Store(ctx, &archivist.StoreRequest{Object: string(objBytes)}); err != nil {
		return err
	}

	return nil
}
