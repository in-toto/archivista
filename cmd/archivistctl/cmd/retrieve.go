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
	"github.com/testifysec/archivist/client"
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

			return client.Download(cmd.Context(), archivistUrl, args[0], out)
		},
	}

	subjectCmd = &cobra.Command{
		Use:          "subjects",
		Short:        "Retrieves all subjects on an in-toto statement by the envelope gitoid",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := client.GraphQlQuery[retrieveSubjectResults](cmd.Context(), archivistUrl, retrieveSubjectsQuery, retrieveSubjectVars{Gitoid: args[0]})
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
	for _, edge := range results.Dsses.Edges {
		for _, subject := range edge.Node.Statement.Subjects {
			digestStrings := make([]string, 0, len(subject.SubjectDigest))
			for _, digest := range subject.SubjectDigest {
				digestStrings = append(digestStrings, fmt.Sprintf("%s:%s", digest.Algorithm, digest.Value))
			}

			fmt.Printf("Name: %s\nDigests: %s\n", subject.Name, strings.Join(digestStrings, ", "))
		}
	}
}

type retrieveSubjectVars struct {
	Gitoid string `json:"gitoid"`
}

type retrieveSubjectResults struct {
	Dsses struct {
		Edges []struct {
			Node struct {
				Statement struct {
					Subjects []struct {
						Name          string `json:"name"`
						SubjectDigest []struct {
							Algorithm string `json:"algorithm"`
							Value     string `json:"value"`
						} `json:"subject_digest"`
					} `json:"subjects"`
				} `json:"statement"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"dsses"`
}

const retrieveSubjectsQuery = `query($gitoid: String!) {
	dsses(
		where: {
			gitoidSha256: $gitoid 
		}
	) {
		edges {
			node {
				statement {
					subjects {
						name
						subject_digest {
							algorithm
							value
						}
					}
				}
			}
		}
	}
}`
