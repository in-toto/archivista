package judgeapi

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	kratos "github.com/ory/kratos-client-go"
	"github.com/sirupsen/logrus"
	"gitlab.com/testifysec/judge-platform/judge-api/ent"
	"gitlab.com/testifysec/judge-platform/judge-api/internal/resolvers"
	"gitlab.com/testifysec/judge-platform/judge-api/viewer"
)

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id uuid.UUID) (ent.Noder, error) {
	panic(fmt.Errorf("not implemented"))
}

// Nodes is the resolver for the nodes field.
func (r *queryResolver) Nodes(ctx context.Context, ids []uuid.UUID) ([]ent.Noder, error) {
	panic(fmt.Errorf("not implemented"))
}

// Projects is the resolver for the projects field.
func (r *queryResolver) Projects(ctx context.Context) ([]*ent.Project, error) {
	viewer := viewer.FromContext(ctx)
	if viewer == nil {
		return nil, fmt.Errorf("viewer not found")
	}

	// Get the identity ID from the viewer
	id := viewer.UserIdentityID()

	api := kratos.NewAPIClient(&kratos.Configuration{
		Host:          "",
		Scheme:        "",
		DefaultHeader: nil,
		UserAgent:     "",
		Debug:         false,
		Servers: kratos.ServerConfigurations{kratos.ServerConfiguration{
			URL:         r.client.Config.KratosAdminUrl,
			Description: "",
			Variables:   nil,
		}},
		OperationServers: nil,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 10 * time.Second,
				}).DialContext,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	})

	identities, _, err := api.IdentityApi.GetIdentity(context.Background(), id).IncludeCredential([]string{"oidc"}).Execute()
	if err != nil {
		logrus.Error("failed to get identity", err)
		return nil, err
	}

	oidcCredsConfig := identities.GetCredentials()["oidc"].Config
	providers, exists := oidcCredsConfig["providers"]
	if !exists {
		logrus.Error("no providers found", err)
		return nil, err
	}

	jsonbody, err := json.Marshal(providers)
	if err != nil {
		return nil, err
	}

	ps := []resolvers.Provider{}
	if err := json.Unmarshal(jsonbody, &ps); err != nil {
		return nil, err
	}

	token := ps[0].InitialAccessToken

	c := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var (
		query   = resolvers.ViewerReposQuery
		perPage = 100
		repos   []resolvers.GhRepoDTO
	)

	for {
		responseBody, err := resolvers.SendGraphQLRequest(ctx, c, query, nil, token)
		if err != nil {
			logrus.WithContext(ctx).WithError(err).Error("failed to get repos from github")
			return nil, err
		}

		reposResponse := resolvers.ReposResponseData{}
		if err := json.Unmarshal(responseBody, &reposResponse); err != nil {
			logrus.WithContext(ctx).WithError(err).Error("failed to unmarshal repos from github")
			return nil, err
		}

		repoList := reposResponse.Data.Viewer.Repositories.Nodes
		repos = append(repos, repoList...)

		if !reposResponse.Data.Viewer.Repositories.PageInfo.HasNextPage {
			break
		}

		query = fmt.Sprintf(`{
			viewer {
				repositories(first: %d, after: "%s", ownerAffiliations:[OWNER, COLLABORATOR, ORGANIZATION_MEMBER]) {
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
		}`, perPage, reposResponse.Data.Viewer.Repositories.PageInfo.EndCursor)
	}

	res := make([]*ent.Project, len(repos))
	for i, repo := range repos {

		res[i] = &ent.Project{
			ID:         uuid.New(),
			RepoID:     repo.ID,
			Name:       repo.NameWithOwner,
			Projecturl: repo.URL,
		}
	}

	return res, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
