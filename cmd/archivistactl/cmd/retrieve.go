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
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/in-toto/archivista/pkg/api"
	"github.com/spf13/cobra"
)

var (
	outFile string

	retrieveCmd = &cobra.Command{
		Use:          "retrieve",
		Short:        "Retrieve information from an archivista server",
		SilenceUsage: true,
	}

	envelopeCmd = &cobra.Command{
		Use:          "envelope",
		Short:        "Retrieves a dsse envelope by it's gitoid from archivista",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var out io.Writer = os.Stdout
			if len(outFile) > 0 {
				file, err := os.Create(outFile)
				if err != nil {
					return err
				}

				defer file.Close()
				out = file
			}
			return api.DownloadWithWriter(cmd.Context(), archivistaUrl, args[0], out)
		},
	}

	subjectCmd = &cobra.Command{
		Use:          "subjects",
		Short:        "Retrieves all subjects on an in-toto statement by the envelope gitoid",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := api.GraphQlQuery[api.RetrieveSubjectResults](
				cmd.Context(),
				archivistaUrl,
				api.RetrieveSubjectsQuery,
				api.RetrieveSubjectVars{Gitoid: args[0]},
			)
			if err != nil {
				return err
			}

			printSubjects(results)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(retrieveCmd)
	retrieveCmd.AddCommand(envelopeCmd)
	retrieveCmd.AddCommand(subjectCmd)
	envelopeCmd.Flags().StringVarP(&outFile, "out", "o", "", "File to write the envelope out to. Defaults to stdout")
}

func printSubjects(results api.RetrieveSubjectResults) {
	for _, edge := range results.Subjects.Edges {
		digestStrings := make([]string, 0, len(edge.Node.SubjectDigests))
		for _, digest := range edge.Node.SubjectDigests {
			digestStrings = append(digestStrings, fmt.Sprintf("%s:%s", digest.Algorithm, digest.Value))
		}

		rootCmd.Printf("Name: %s\nDigests: %s\n", edge.Node.Name, strings.Join(digestStrings, ", "))
	}
}
