package auth

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

	type MetadataPublic struct {
		AssignedTenants []uuid.UUID `json:"assigned_tenants"`
	}

	var metadata MetadataPublic

	metaDataBytes, ok := metaDataPublic.([]byte)
	if !ok {
		logrus.Errorf("Failed to convert metaDataPublic to []byte")
		return v, fmt.Errorf("metaDataPublic conversion error")
	}

	err = json.Unmarshal(metaDataBytes, &metadata)
	if err != nil {
		logrus.Errorf("Failed to unmarshal metadata: %v", err)
		return v, err
	}

	v.TenantsAssigned = metadata.AssignedTenants

	return v, nil
}
