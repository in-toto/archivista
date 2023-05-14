package resolvers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

const (
	ViewerReposQuery = `{
		viewer {
			repositories(first: 100, ownerAffiliations:[OWNER, COLLABORATOR, ORGANIZATION_MEMBER]) {
				totalCount
				nodes {
				  id
				  nameWithOwner
				  isPrivate
				  isInOrganization
				  url
				}
				pageInfo {
					endCursor
					hasNextPage
				}
			}
		}
	}`
)

type Provider struct {
	InitialIDToken      string `json:"initial_id_token"`
	Subject             string `json:"subject"`
	Provider            string `json:"provider"`
	InitialAccessToken  string `json:"initial_access_token"`
	InitialRefreshToken string `json:"initial_refresh_token"`
}
type GhRepoDTO struct {
	ID               string `json:"id"`
	NameWithOwner    string `json:"nameWithOwner"`
	IsInOrganization bool   `json:"isInOrganization"`
	IsPrivate        bool   `json:"isPrivate"`
	LatestRelease    struct {
		ID string `json:"id"`
	} `json:"latestRelease"`
	URL string `json:"url"`
}
type ReposResponseData struct {
	Data struct {
		Viewer struct {
			Repositories struct {
				TotalCount int         `json:"totalCount"`
				Nodes      []GhRepoDTO `json:"nodes"`
				PageInfo   struct {
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
				} `json:"pageInfo"`
			} `json:"repositories"`
		} `json:"viewer"`
	} `json:"data"`
}

func SendGraphQLRequest(ctx context.Context, client *http.Client, query string, variables map[string]interface{}, token string) ([]byte, error) {
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.github.com/graphql", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	spew.Dump(string(responseBody))

	return responseBody, nil
}
