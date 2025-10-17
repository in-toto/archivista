package vaultauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	authv1 "k8s.io/api/authentication/v1"

	"github.com/davepgreene/go-db-credential-refresh/store/vault/vaulttest"
)

const (
	testCACert = `-----BEGIN CERTIFICATE-----
MIIC5zCCAc+gAwIBAgIBATANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwptaW5p
a3ViZUNBMB4XDTE5MDEwNTE4MDkxNFoXDTI5MDEwMzE4MDkxNFowFTETMBEGA1UE
AxMKbWluaWt1YmVDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMNw
BYLyqYEVm4vbik0NibQ7+G414dPUUZc3UCScrvMYASa+Krcc8J8Ic0TeDdsluiYs
hujALbu+LtFNYeIpMBgZPUaBVOtSrnBe9ieG0XZmxDa303uz2awzYkivWab58Tsx
RLojX2z4ZJUXhb1m6VN96x07tf4MujnQgmfm3GZ/cMn/BUaTSJOKXiKTDTys6dbz
U3UyvQnxP9QkWloU6HICqPObzpY/kkLdsOWPfiGn2lINZ/9zkeW8Qe9QalKRuGI2
2+ZWOTZyREvfln/3LML8q9kAmk54NMtSG3mGCgDOL+HsRVNnqZC2QBHmHJGD+2nz
z6C1iSV0W4ZgDR0HSIMCAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgKkMB0GA1UdJQQW
MBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4IBAQC1epomLINu4JVBCZpaHDpWg5cAHyxAVggKihCP+hIsOkfsmyTa
RPCLNASmW/PbDDzJKAzQaC23KDFW9WCqr1MlgsJhZMW8tkiPegL18DGxupjwzIIM
meiZoPBEpFGz0JVhfu0FMIVbvKjhuBgbZd3rKZEFHZMer7L+ZZ2Pd/5UY1s7oslq
Z938fecvWwIQLHE+Jar1KqvdlLlP798+w7G5de3gIN0svflcpbd8+w3X7h+dzouu
qevae2NSZJ5r8Fo5Ch3sI63c6GCoUaMM5Ho7mHUM32BeGxy99Z3G6364akR3I819
qQYZl8EZf4Jznaes/XFP0Yb+IhGXBoR9Ib+I
-----END CERTIFICATE-----`

	jwtData     = `eyJhbGciOiJSUzI1NiIsImtpZCI6IlpJOEY4RHVoMktrY0JxTjhGSGxyMEhER2l2OEtFR2xFSnlITUZRc1UwZ28ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InZhdWx0LWF1dGgtdG9rZW4tdmQ0bjQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoidmF1bHQtYXV0aCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjlhZjM3NjRlLWZmZDMtNDJiZC1hZjVkLTE2MzUwZTM0NjkyYyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OnZhdWx0LWF1dGgifQ.ZfkKFqeAIaNXmk-i7LwrXXoOIjv4WlQ1gHFOXHpSo0Wdq16KKu1VOnCkzUh9bIApL5pIXZu4-eYwP2SwokRafXBY_5znqvXoI3F1fxmw25jBT9ZeyDEKZOxyO7mtHnh7LZQ_pBUPPflClhAwacbBrTjnIpHoiWq-Z1_BeuenlRdBYQYjdXEOPK-W1bFbCqx4hq_x91v-JMAcJqQUf0ZSY3jU-vcAOmFfv_0S4K2_syUyfkYVPr_pX-0wOvwkv0nDhV-fhqux51onQyYDd_gejvjGvviDJcbXxT4sIYgbS8IKtRwI3lAhpQQyuaQbVI6DKASs9z-jvvg0VO7T2FMFIw`
	jwtUID      = "9af3764e-ffd3-42bd-af5d-16350e34692c"
	jwtUsername = "system:serviceaccount:default:vault-auth"
)

var jwtGroups = []string{
	"system:serviceaccounts",
	"system:serviceaccounts:default",
	"system:authenticated",
}

func getDockerHostIP() string {
	switch runtime.GOOS {
	case "darwin":
		return "host.docker.internal"
	case "windows":
		return "host.docker.internal"
	case "linux":
		return "172.17.0.1"
	default:
		return ""
	}
}

func runningInGithubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

// this HandlerFunc mocks out an response from k8s's tokenreviews endpoint.
func tokenReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/apis/authentication.k8s.io/v1/tokenreviews" {
		w.WriteHeader(404)

		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)

		return
	}
	defer r.Body.Close() //nolint:errcheck

	var tr authv1.TokenReview
	if err := json.Unmarshal(body, &tr); err != nil {
		w.WriteHeader(500)

		return
	}

	tr.Status = authv1.TokenReviewStatus{
		Authenticated: true,
		User: authv1.UserInfo{
			UID:      jwtUID,
			Username: jwtUsername,
			Groups:   jwtGroups,
		},
	}

	json.NewEncoder(w).Encode(tr) //nolint:errcheck
}

func TestKubernetesAuth(t *testing.T) {
	if runningInGithubActions() {
		t.Skip("Skipping until testcontainers-go supports docker networking and network aliases")
	}

	srv := &httptest.Server{
		Listener: func() net.Listener {
			// We need to bind to 0.0.0.0 to ensure that the docker container can hit this test server
			l, err := net.Listen("tcp4", "0.0.0.0:0")
			if err != nil {
				panic(fmt.Sprintf("httptest: failed to listen on a port: %v", err))
			}

			return l
		}(),
		Config: &http.Server{Handler: http.HandlerFunc(tokenReviewHandler)},
	}

	srv.Start()
	defer srv.Close()

	ctx := context.Background()

	tokenAndClient, vaultContainer, err := vaulttest.CreateTestVault(ctx)
	if err != nil {
		if vaultContainer != nil {
			if err := vaultContainer.Terminate(ctx); err != nil {
				t.Fatal(err)
			}
		}
		t.Fatal(err)
	}
	defer func() {
		if err := vaultContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	serverUrl, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	client := tokenAndClient.Client

	if _, err := client.Auth.KubernetesConfigureAuth(ctx, schema.KubernetesConfigureAuthRequest{
		KubernetesHost:   fmt.Sprintf("http://%s:%s", getDockerHostIP(), serverUrl.Port()),
		KubernetesCaCert: testCACert,
	}); err != nil {
		t.Fatal(err)
	}

	role := "example"
	userName := strings.Split(jwtUsername, ":")

	if _, err := client.Auth.KubernetesWriteAuthRole(ctx, role, schema.KubernetesWriteAuthRoleRequest{
		BoundServiceAccountNames:      []string{userName[len(userName)-1]},
		BoundServiceAccountNamespaces: []string{"default"},
	}); err != nil {
		t.Fatal(err)
	}

	tmpfile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) //nolint:errcheck

	if _, err = tmpfile.Write([]byte(jwtData)); err != nil {
		t.Fatal(err)
	}
	if err = tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	a := NewKubernetesAuth(role, tmpfile.Name())
	token, err := a.GetToken(ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatal("expected a token but didn't get one")
	}

	// Verify the token
	if err := client.SetToken(token); err != nil {
		t.Fatal(err)
	}

	resp, err := client.Auth.TokenLookUpSelf(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("expected a valid response")
	}

	if resp.Data == nil {
		t.Fatal("expected response to have data")
	}

	path, ok := resp.Data["path"]
	if !ok {
		t.Fatal("expected 'path' to be in auth response data")
	}

	if path != "auth/kubernetes/login" {
		t.Fatalf("expected 'path' to be k8s login path but got %s instead", path)
	}
}

func TestKubernetesAuthFileError(t *testing.T) {
	p := "/foo/bar/baz"
	k := NewKubernetesAuth("role", p)
	client, err := vault.New()
	if err != nil {
		t.Fatal(err)
	}
	_, err = k.GetToken(context.Background(), client)
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	pathErr := &os.PathError{
		Op:   "open",
		Path: p,
		Err:  errors.New("no such file or directory"),
	}
	if err.Error() != pathErr.Error() {
		t.Fatalf("expected error to be '%v' but got '%v' instead", pathErr, err)
	}
}

func TestKubernetesAuthVaultError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer srv.Close()

	client, err := vault.New(vault.WithAddress(srv.URL))
	if err != nil {
		t.Fatal(err)
	}

	tmpfile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) //nolint:errcheck

	if _, err = tmpfile.Write([]byte(jwtData)); err != nil {
		t.Fatal(err)
	}
	if err = tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	k := NewKubernetesAuth("role", tmpfile.Name())
	_, err = k.GetToken(context.Background(), client)
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	var respErr *vault.ResponseError
	if errors.As(err, &respErr) {
		if respErr.StatusCode != http.StatusNotFound {
			t.Fatalf(
				"expected to get a %d but got a %d instead",
				http.StatusNotFound,
				respErr.StatusCode,
			)
		}

		loginURL := fmt.Sprintf("%s/v1/auth/kubernetes/login", srv.URL)
		responseErrorURL := respErr.OriginalRequest.URL.String()
		if responseErrorURL != loginURL {
			t.Fatalf("expected URL to be %s but got %s instead", loginURL, responseErrorURL)
		}
		if respErr.OriginalRequest.Method != http.MethodPost {
			t.Fatalf(
				"expected method %s but got %s instead",
				http.MethodPut,
				respErr.OriginalRequest.Method,
			)
		}
	}
}
