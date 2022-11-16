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
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	archivistapi "github.com/testifysec/archivist-api"
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
			var out io.Writer = os.Stdout
			if len(outFile) > 0 {
				file, err := os.Create(outFile)
				if err != nil {
					return err
				}

				defer file.Close()
				out = file
			}

			return archivistapi.DownloadWithWriter(cmd.Context(), archivistUrl, args[0], out)
		},
	}

	subjectCmd = &cobra.Command{
		Use:          "subjects",
		Short:        "Retrieves all subjects on an in-toto statement by the envelope gitoid",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := archivistapi.GraphQlQuery[retrieveSubjectResults](cmd.Context(), archivistUrl, retrieveSubjectsQuery, retrieveSubjectVars{Gitoid: args[0]})
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

func printSubjects(results retrieveSubjectResults) {
	for _, edge := range results.Subjects.Edges {
		digestStrings := make([]string, 0, len(edge.Node.SubjectDigests))
		for _, digest := range edge.Node.SubjectDigests {
			digestStrings = append(digestStrings, fmt.Sprintf("%s:%s", digest.Algorithm, digest.Value))
		}

		fmt.Printf("Name: %s\nDigests: %s\n", edge.Node.Name, strings.Join(digestStrings, ", "))
	}
}

type retrieveSubjectVars struct {
	Gitoid string `json:"gitoid"`
}

type retrieveSubjectResults struct {
	Subjects struct {
		Edges []struct {
			Node struct {
				Name           string `json:"name"`
				SubjectDigests []struct {
					Algorithm string `json:"algorithm"`
					Value     string `json:"value"`
				} `json:"subjectDigests"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"subjects"`
}

const retrieveSubjectsQuery = `query($gitoid: String!) {
	subjects(
		where: {
			hasStatementWith:{
        hasDsseWith:{
          gitoidSha256: $gitoid
        }
      } 
		}
	) {
		edges {
      node{
        name
        subjectDigests{
          algorithm
          value
        }
      }
    }
  }
}`
