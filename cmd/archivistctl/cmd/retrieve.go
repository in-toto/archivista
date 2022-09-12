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
	"bytes"
	"context"
	"fmt"
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
		Short:        "Retrieve information from an archivist server",
		SilenceUsage: true,
	}

	envelopeCmd = &cobra.Command{
		Use:          "envelope",
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
				if err != nil {
					return err
				}

				defer file.Close()
				out = file
			}

			return retrieveEnvelope(cmd.Context(), archivist.NewCollectorClient(conn), args[0], out)
		},
	}

	subjectCmd = &cobra.Command{
		Use:          "subjects",
		Short:        "Retrieves all subjects on an in-toto statement by the envelope gitoid",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := newConn(archivistUrl)
			if err != nil {
				return err
			}

			return retrieveSubjects(cmd.Context(), archivist.NewArchivistClient(conn), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(retrieveCmd)
	retrieveCmd.AddCommand(envelopeCmd)
	retrieveCmd.AddCommand(subjectCmd)
	envelopeCmd.Flags().StringVarP(&outFile, "out", "o", "", "File to write the envelope out to. Defaults to stdout")
}

func retrieveSubjects(ctx context.Context, client archivist.ArchivistClient, gitoid string) error {
	stream, err := client.GetSubjects(ctx, &archivist.GetSubjectsRequest{EnvelopeGitoid: gitoid})
	if err != nil {
		return err
	}

	for {
		subject, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		fmt.Printf("Name: %s\nDigests:\n%s\n", subject.GetName(), digestString(subject.GetDigest()))
	}

	return nil
}

func digestString(digest map[string]string) string {
	sb := strings.Builder{}
	for algo, value := range digest {
		sb.WriteString(fmt.Sprintf("Algo: %s\nValue: %s\n", algo, value))
	}

	return sb.String()
}

func retrieveEnvelope(ctx context.Context, client archivist.CollectorClient, gitoid string, out io.Writer) error {
	stream, err := client.Get(ctx, &archivist.GetRequest{Gitoid: gitoid})
	if err != nil {
		return err
	}

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if _, err := io.Copy(out, bytes.NewReader(chunk.GetChunk())); err != nil {
			return err
		}
	}

	return nil
}
