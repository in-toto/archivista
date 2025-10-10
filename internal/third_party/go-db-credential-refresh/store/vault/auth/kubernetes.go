package vaultauth

import (
	"context"
	"os"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

const (
	kubeTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token" //nolint:gosec
)

// KubernetesAuth gets the vault auth token from the kubernetes secrets file.
type KubernetesAuth struct {
	role string
	path string
}

// NewKubernetesAuth creates a new k8s secret auth token location.
func NewKubernetesAuth(role, path string) *KubernetesAuth {
	if path == "" {
		path = kubeTokenPath
	}

	return &KubernetesAuth{
		role: role,
		path: path,
	}
}

// GetToken implements the TokenLocation interface.
func (k *KubernetesAuth) GetToken(ctx context.Context, client *vault.Client) (string, error) {
	token, err := os.ReadFile(k.path)
	if err != nil {
		return "", err
	}

	secret, err := client.Auth.KubernetesLogin(ctx, schema.KubernetesLoginRequest{
		Jwt:  string(token),
		Role: k.role,
	})
	if err != nil {
		return "", err
	}

	return secret.Auth.ClientToken, nil
}
