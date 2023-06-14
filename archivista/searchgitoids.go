// Copyright 2022 The Witness Contributors
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

package archivista

import (
	"context"

	archivistaapi "github.com/testifysec/archivista-api"
)

type searchGitoidResponse struct {
	Dsses struct {
		Edges []struct {
			Node struct {
				Gitoid string `json:"gitoidSha256"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"dsses"`
}

type SearchGitoidVariables struct {
	SubjectDigests []string `json:"subjectDigests"`
	CollectionName string   `json:"collectionName"`
	Attestations   []string `json:"attestations"`
	ExcludeGitoids []string `json:"excludeGitoids"`
}

func (c *Client) SearchGitoids(ctx context.Context, vars SearchGitoidVariables) ([]string, error) {
	const query = `query ($subjectDigests: [String!], $attestations: [String!], $collectionName: String!, $excludeGitoids: [String!]) {
  dsses(
    where: {
			gitoidSha256NotIn: $excludeGitoids,
			hasStatementWith: {
				hasAttestationCollectionsWith: {
					name: $collectionName,
					hasAttestationsWith: {
						typeIn: $attestations
					}
				},
				hasSubjectsWith: {
					hasSubjectDigestsWith: {
						valueIn: $subjectDigests
					}
				}
			}
		}
  ) {
    edges {
      node {
        gitoidSha256
      }
    }
  }
}`

	response, err := archivistaapi.GraphQlQuery[searchGitoidResponse](ctx, c.url, query, vars)
	if err != nil {
		return nil, err
	}

	gitoids := make([]string, 0, len(response.Dsses.Edges))
	for _, edge := range response.Dsses.Edges {
		gitoids = append(gitoids, edge.Node.Gitoid)
	}

	return gitoids, nil
}
