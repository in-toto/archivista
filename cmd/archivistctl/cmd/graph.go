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
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/spf13/cobra"
	archivistapi "github.com/testifysec/archivist-api"
)

var (
	backRefs = []string{
		"https://witness.dev/attestations/git/v0.1/commithash",
		"https://witness.dev/attestations/gitlab/v0.1/pipelineurl",
	}

	graphIterations = 1
	graphCmd        = &cobra.Command{
		Use:          "graph",
		Short:        "Create a graphviz file representing the attestation graph",
		SilenceUsage: true,
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
			_, digest, err := validateDigestString(args[0])
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

			g := graphviz.New()
			graph, err := g.Graph()
			if err != nil {
				return err
			}

			defer g.Close()
			digests := []string{digest}
			envelopesByDigest := make(map[string]map[string]struct{})
			for i := 0; i <= graphIterations; i++ {
				results, err := archivistapi.GraphQlQuery[graphResults](cmd.Context(), archivistUrl, graphQuery, graphVars{Digests: digests})
				if err != nil {
					return err
				}

				newDigests, err := addNodesToGraph(graph, results, envelopesByDigest)
				if err != nil {
					return nil
				}

				digests = newDigests
			}

			for digest, envelopes := range envelopesByDigest {
				for envelope := range envelopes {
					start, err := graph.Node(envelope)
					if err != nil {
						return err
					}

					for endEnv := range envelopes {
						if endEnv == envelope {
							continue
						}

						end, err := graph.Node(endEnv)
						if err != nil {
							return err
						}

						if _, err := graph.CreateEdge(digest, start, end); err != nil {
							return err
						}
					}
				}
			}

			return g.Render(graph, graphviz.XDOT, out)
		},
	}
)

func init() {
	graphCmd.Flags().StringVarP(&outFile, "out", "o", "", "File to write the graph to. Defaults to stdout")
	rootCmd.AddCommand(graphCmd)
}

func addNodesToGraph(g *cgraph.Graph, results graphResults, envelopesByDigest map[string]map[string]struct{}) ([]string, error) {
	digests := make(map[string]struct{}, 0)
	for _, dsse := range results.Dsses.Edges {
		n, err := g.CreateNode(dsse.Node.GitoidSha256)
		if err != nil {
			return []string{}, err
		}

		n.SetLabel(fmt.Sprintf("%v\n%v", dsse.Node.Statement.AttestationCollection.Name, dsse.Node.GitoidSha256))

		for _, subject := range dsse.Node.Statement.Subjects.Edges {
			for _, backRef := range backRefs {
				for _, digest := range subject.Node.SubjectDigests {
					if strings.HasPrefix(subject.Node.Name, backRef) {
						digests[digest.Value] = struct{}{}
					}

					envelopes, ok := envelopesByDigest[digest.Value]
					if !ok {
						envelopes = make(map[string]struct{})
					}

					envelopes[dsse.Node.GitoidSha256] = struct{}{}
					envelopesByDigest[digest.Value] = envelopes
				}
			}
		}
	}

	newDigests := make([]string, 0)
	for digest := range digests {
		newDigests = append(newDigests, digest)
	}

	return newDigests, nil
}

type graphVars struct {
	Digests []string `json:"digests"`
}

type graphResults struct {
	Dsses struct {
		Edges []struct {
			Node struct {
				GitoidSha256 string `json:"gitoidSha256"`
				Statement    struct {
					AttestationCollection struct {
						Name string `json:"name"`
					} `json:"attestationCollections"`
					Subjects struct {
						Edges []struct {
							Node struct {
								Name           string `json:"name"`
								SubjectDigests []struct {
									Value     string `json:"value"`
									Algorithm string `json:"algorithm"`
								} `json:"subjectDigests"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"subjects"`
				} `json:"statement"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"dsses"`
}

const graphQuery = `query($digests: [String!]) {
  dsses(
    where: {
      hasStatementWith: {
        hasSubjectsWith: {
          hasSubjectDigestsWith: {
            valueIn: $digests
          }
        }
      }
    }
  ) {
    totalCount
    edges {
      node {
        gitoidSha256
        statement {
					attestationCollections {
						name
					}
          subjects(where:{not:{nameHasPrefix:"https://witness.dev/attestations/product/v0.1"}}){
            edges{
              node{
                name
                subjectDigests{
                  value
                  algorithm
                }
              }
            }
          }
        }
      }
    }
  }
}`
