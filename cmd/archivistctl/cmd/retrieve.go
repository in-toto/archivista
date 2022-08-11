package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
			conn, err := newConn(archivistGrpcUrl)
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

	subjectCmd = &cobra.Command{
		Use:          "subjects",
		Short:        "Retrieves all subjects on an in-toto statement by the envelope gitoid",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return retrieveSubjects(cmd.Context(), archivistGqlUrl, args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(retrieveCmd)
	retrieveCmd.AddCommand(envelopeCmd)
	retrieveCmd.AddCommand(subjectCmd)
	envelopeCmd.Flags().StringVarP(&outFile, "out", "o", "", "File to write the envelope out to. Defaults to stdout")
}

func retrieveSubjects(ctx context.Context, graphUrl, gitoid string) error {
	query, err := json.Marshal(subjectQuery(gitoid))
	if err != nil {
		return err
	}

	results, err := executeGraphQlQuery(graphUrl, query)
	if err != nil {
		return err
	}

	parsedResults := struct {
		Data struct {
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
		} `json:"data"`
	}{}

	if err := json.Unmarshal(results, &parsedResults); err != nil {
		return err
	}

	for _, edge := range parsedResults.Data.Dsses.Edges {
		for _, subject := range edge.Node.Statement.Subjects {
			digestStrings := make([]string, 0, len(subject.SubjectDigest))
			for _, digest := range subject.SubjectDigest {
				digestStrings = append(digestStrings, fmt.Sprintf("%s:%s", digest.Algorithm, digest.Value))
			}

			fmt.Printf("Name: %s\nDigests: %s\n", subject.Name, strings.Join(digestStrings, ", "))
		}
	}

	return nil
}

type digest struct {
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

func executeGraphQlQuery(url string, query []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(query))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	results, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func subjectQuery(gitoid string) map[string]string {
	return map[string]string{
		"query": fmt.Sprintf(`{
	dsses(
		where: {
			gitoidSha256: "%v"
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
}`, gitoid),
	}
}
