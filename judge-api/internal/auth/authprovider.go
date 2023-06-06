package auth

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"

	kratos "github.com/ory/kratos-client-go"
	"github.com/sirupsen/logrus"
	"gitlab.com/testifysec/judge-platform/judge-api/viewer"
)

type AuthProvider interface {
	ValidateAndGetViewer(ctx context.Context, c string) (viewer.UserViewer, error)
}

type KratosAuthProvider struct {
	kratosClient *kratos.APIClient
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

	return &KratosAuthProvider{
		kratosClient: api,
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