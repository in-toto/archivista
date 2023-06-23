package auth

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/networkservicemesh/sdk/pkg/tools/log"
	"net/http"

	kratos "github.com/ory/kratos-client-go"
	"github.com/sirupsen/logrus"
	"gitlab.com/testifysec/judge-platform/judge-api/viewer"
)

type AuthProvider interface {
	ValidateAndGetViewer(ctx context.Context, c string) (viewer.UserViewer, error)
}

type KratosAuthProvider struct {
	kratosClient      *kratos.APIClient
	kratosAdminClient *kratos.APIClient
}

type MetadataPublic struct {
	AssignedTenants []string `json:"assigned_tenants"`
}

func NewKratosAuthProvider() *KratosAuthProvider {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	api := kratos.NewAPIClient(&kratos.Configuration{
		Host:       "kratos-public.kratos.svc.cluster.local",
		Scheme:     "http",
		Debug:      true,
		Servers:    []kratos.ServerConfiguration{{URL: "http://kratos-public.kratos.svc.cluster.local"}},
		HTTPClient: client,
	})

	adminApi := kratos.NewAPIClient(&kratos.Configuration{
		Host:       "kratos-admin.kratos.svc.cluster.local",
		Scheme:     "http",
		Debug:      true,
		Servers:    []kratos.ServerConfiguration{{URL: "http://kratos-admin.kratos.svc.cluster.local"}},
		HTTPClient: client,
	})

	return &KratosAuthProvider{
		kratosClient:      api,
		kratosAdminClient: adminApi,
	}
}

func (k *KratosAuthProvider) ValidateAndGetViewer(ctx context.Context, c string) (viewer.UserViewer, error) {
	v := viewer.UserViewer{}
	api := k.kratosClient

	r := api.FrontendApi.ToSession(ctx).Cookie(c)
	session, _, err := r.Execute()
	if err != nil {
		logrus.Errorf("Failed to get session: %v", err)
		return v, err
	}

	v.IdentityID = session.Identity.Id

	metaDataPublic := session.Identity.GetMetadataPublic()

	var metadata MetadataPublic

	m, err := json.Marshal(metaDataPublic)
	if err != nil {
		logrus.Errorf("Failed to marshal metadata: %v", err)
		return v, err
	}
	err = json.Unmarshal(m, &metadata)
	if err != nil {
		logrus.Errorf("Failed to unmarshal metadata: %v", err)
		return v, err
	}

	for _, tenant := range metadata.AssignedTenants {
		parsedTenant, _ := uuid.Parse(tenant)
		v.TenantsAssigned = append(v.TenantsAssigned, parsedTenant)
	}

	return v, nil
}

func (k *KratosAuthProvider) UpdateAssignedTenantsWithIdentityId(w http.ResponseWriter, r *http.Request) {
	type MetadataWebhookRequest struct {
		IdentityId string                 `json:"identityId"`
		Traits     map[string]interface{} `json:"traits"`
	}

	var metadataWebhookRequest MetadataWebhookRequest
	err := json.NewDecoder(r.Body).Decode(&metadataWebhookRequest)
	if err != nil {
		log.FromContext(r.Context()).Errorf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading webhook request body"))
		return
	}

	var metadata MetadataPublic
	metadata.AssignedTenants = []string{metadataWebhookRequest.IdentityId}
	_, _, err = k.kratosAdminClient.IdentityApi.UpdateIdentity(r.Context(), metadataWebhookRequest.IdentityId).UpdateIdentityBody(kratos.UpdateIdentityBody{
		Traits:         metadataWebhookRequest.Traits,
		MetadataPublic: metadata,
	}).Execute()
	if err != nil {
		log.FromContext(r.Context()).Errorf("Error updating identity: %v", err)
	}
}
