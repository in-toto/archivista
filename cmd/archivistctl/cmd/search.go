package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
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
			algo, digest, err := validateDigestString(args[0])
			if err != nil {
				return err
			}

			return searchArchivist(archivistGqlUrl, algo, digest)
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

func searchArchivist(url, algo, digest string) error {
	query, err := json.Marshal(searchQuery(algo, digest))
	if err != nil {
		return nil
	}

	results, err := executeGraphQlQuery(url, query)
	if err != nil {
		return err
	}

	parsedResults := struct {
		Data struct {
			Dsses struct {
				Edges []struct {
					Node struct {
						GitoidSha256 string `json:"gitoid_sha256"`
						Statement    struct {
							AttestationCollection struct {
								Name         string `json:"name"`
								Attestations []struct {
									Type string `json:"type"`
								} `json:"attestations"`
							} `json:"attestation_collection"`
						} `json:"statement"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"dsses"`
		} `json:"data"`
	}{}

	if err := json.Unmarshal(results, &parsedResults); err != nil {
		return err
	}

	for _, edge := range parsedResults.Data.Dsses.Edges {
		fmt.Printf("Gitoid: %s\n", edge.Node.GitoidSha256)
		fmt.Printf("Collection name: %s\n", edge.Node.Statement.AttestationCollection.Name)
		types := make([]string, 0, len(edge.Node.Statement.AttestationCollection.Attestations))
		for _, attestation := range edge.Node.Statement.AttestationCollection.Attestations {
			types = append(types, attestation.Type)
		}

		fmt.Printf("Attestations: %s\n\n", strings.Join(types, ", "))
	}

	return nil
}

func searchQuery(algo, value string) map[string]string {
	return map[string]string{
		"query": fmt.Sprintf(`{
  dsses(
    where: {
      hasStatementWith: {
        hasSubjectsWith: {
          hasSubjectDigestsWith: {
            value: "%v", 
            algorithm: "%v"
          }
        }
      }
    }
  ) {
    edges {
      node {
        gitoid_sha256
        statement {
          attestation_collection {
            name
            attestations {
              type
            }
          }
        }
      }
    }
  }
}`, value, algo),
	}
}
