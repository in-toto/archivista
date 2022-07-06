package cmd

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/testifysec/archivist-api/pkg/api/archivist"
)

var (
	outFile string

	retrieveCmd = &cobra.Command{
		Use:          "retrieve",
		Short:        "Retrieves a dsse envelope by it's gitoid from archivist",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := newConn(archivistUrl)
			if err != nil {
				return err
			}

			var out io.Writer = os.Stdout
			if len(outFile) > 0 {
				file, err := os.Create(outFile)
				defer file.Close()
				if err != nil {
					return err
				}

				out = file
			}

			return retrieveEnvelope(cmd.Context(), archivist.NewCollectorClient(conn), args[0], out)
		},
	}
)

func init() {
	rootCmd.AddCommand(retrieveCmd)
	retrieveCmd.Flags().StringVarP(&outFile, "out", "o", "", "File to write the envelope out to. Defaults to stdout")
}

func retrieveEnvelope(ctx context.Context, client archivist.CollectorClient, gitoid string, out io.Writer) error {
	resp, err := client.Get(ctx, &archivist.GetRequest{Gitoid: gitoid})
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, strings.NewReader(resp.Object)); err != nil {
		return err
	}

	return nil
}
