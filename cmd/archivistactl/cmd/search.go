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
	"errors"
	"strings"

	"github.com/in-toto/archivista/pkg/api"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:          "search",
	Short:        "Searches the archivista instance for an attestation matching a query",
	SilenceUsage: true,
	Long: `Searches the archivista instance for an envelope with a specified subject digest.
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
		algo, digest, err := validateDigestString(args[0])
		if err != nil {
			return err
		}

		results, err := api.GraphQlQuery[api.SearchResults](
			cmd.Context(),
			archivistaUrl,
			api.SearchQuery,
			api.SearchVars{Algorithm: algo, Digest: digest},
		)
		if err != nil {
			return err
		}

		// TODO(nick): e2e test are failing here because of output, this is kind of terrible
		printResults(results)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func validateDigestString(ds string) (algo, digest string, err error) {
	algo, digest, found := strings.Cut(ds, ":")
	if !found {
		return "", "", errors.New("invalid digest string. expected algorithm:digest")
	}

	return algo, digest, nil
}

func printResults(results api.SearchResults) {
	for _, edge := range results.Dsses.Edges {
		rootCmd.Printf("Gitoid: %s\n", edge.Node.GitoidSha256)
		rootCmd.Printf("Collection name: %s\n", edge.Node.Statement.AttestationCollection.Name)
		types := make([]string, 0, len(edge.Node.Statement.AttestationCollection.Attestations))
		for _, attestation := range edge.Node.Statement.AttestationCollection.Attestations {
			types = append(types, attestation.Type)
		}

		rootCmd.Printf("Attestations: %s\n\n", strings.Join(types, ", "))
	}
}
