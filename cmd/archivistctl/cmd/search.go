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
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/testifysec/archivist-api/pkg/api/archivist"
)

var (
	collectionName string

	searchCmd = &cobra.Command{
		Use:          "search",
		Short:        "Searches the archivist instance for an attestation matching a query",
		SilenceUsage: true,
		Long: `Searches the archivist instance for an envelope with a specified subject digest.
Optionally a collection name can be provided to further constrain results.

Digests are expected to be in the form algorithm:digest, for instance: sha256:456c0c9a7c05e2a7f84c139bbacedbe3e8e88f9c`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("expected exactly 1 argument")
			}

			if _, _, err := validateDigestString(args[0]); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := newConn(archivistUrl)
			if err != nil {
				return nil
			}

			defer conn.Close()
			algo, digest, err := validateDigestString(args[0])
			if err != nil {
				return err
			}

			return searchArchivist(cmd.Context(), archivist.NewArchivistClient(conn), algo, digest, collectionName)
		},
	}
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVarP(&collectionName, "collectionname", "n", "", "Only envelopes containing an attestation collection with the provided name will be returned.")
}

func validateDigestString(ds string) (algo, digest string, err error) {
	algo, digest, found := strings.Cut(ds, ":")
	if !found {
		return "", "", errors.New("invalid digest string. expected algorithm:digest")
	}

	return algo, digest, nil
}

func searchArchivist(ctx context.Context, client archivist.ArchivistClient, algo, digest, collName string) error {
	req := &archivist.GetBySubjectDigestRequest{Algorithm: algo, Value: digest, CollectionName: collName}
	stream, err := client.GetBySubjectDigest(ctx, req)
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		printResponse(resp)
	}

	return nil
}

func printResponse(resp *archivist.GetBySubjectDigestResponse) {
	fmt.Printf("Gitoid: %s\nCollection name: %s\nAttestations: %s\n\n", resp.GetGitoid(), resp.GetCollectionName(), strings.Join(resp.GetAttestations(), ", "))
}
