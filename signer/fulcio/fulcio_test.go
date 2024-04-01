// Copyright 2023 The Witness Contributors
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

package fulcio

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"math/big"

	"path/filepath"

	fulciopb "github.com/sigstore/fulcio/pkg/generated/protobuf"
	"github.com/stretchr/testify/require"
	"go.step.sm/crypto/jose"

	"github.com/go-jose/go-jose/v3/jwt"
	"google.golang.org/grpc"
)

func setupFulcioTestService(t *testing.T) (*dummyCAClientService, string) {
	service := &dummyCAClientService{}
	service.server = grpc.NewServer()
	fulciopb.RegisterCAServer(service.server, service)
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	client, err := newClient(context.Background(), "https://localhost", lis.Addr().(*net.TCPAddr).Port, true)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	service.client = client
	go func() {
		if err := service.server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return service, fmt.Sprintf("localhost:%d", lis.Addr().(*net.TCPAddr).Port)
}

func TestNewClient(t *testing.T) {
	ctx := context.Background()

	// test when fulcioURL is empty
	_, err := newClient(ctx, "", 0, false)
	require.Error(t, err)

	// test when fulcioURL is invalid
	_, err = newClient(ctx, "://", 0, false)
	require.Error(t, err)

	// test when connection to Fulcio succeeds
	client, err := newClient(ctx, "https://fulcio.url", 0, false)
	require.NoError(t, err)
	require.NotNil(t, client)
}

type dummyCAClientService struct {
	client fulciopb.CAClient
	server *grpc.Server
	fulciopb.UnimplementedCAServer
}

func (s *dummyCAClientService) GetTrustBundle(ctx context.Context, in *fulciopb.GetTrustBundleRequest) (*fulciopb.TrustBundle, error) {
	return &fulciopb.TrustBundle{
		Chains: []*fulciopb.CertificateChain{},
	}, nil
}

func (s *dummyCAClientService) CreateSigningCertificate(ctx context.Context, in *fulciopb.CreateSigningCertificateRequest) (*fulciopb.SigningCertificate, error) {
	t := &testing.T{}

	cert := fulciopb.SigningCertificate{
		Certificate: &fulciopb.SigningCertificate_SignedCertificateEmbeddedSct{
			SignedCertificateEmbeddedSct: &fulciopb.SigningCertificateEmbeddedSCT{
				Chain: &fulciopb.CertificateChain{
					Certificates: generateCertChain(t),
				},
			},
		},
	}
	return &cert, nil
}
func generateTestToken(email string, subject string) string {

	var claims struct {
		jwt.Claims
		Email   string `json:"email"`
		Subject string `json:"sub"`
	}

	key := []byte("test-secret")
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: key}, nil)
	if err != nil {
		log.Fatal(err)
	}

	if email != "" {
		claims.Email = email
	}

	if subject != "" {
		claims.Subject = subject
	}

	claims.Audience = []string{"sigstore"}

	builder := jwt.Signed(signer).Claims(claims)
	signedToken, _ := builder.CompactSerialize()

	return signedToken
}

func TestGetCert(t *testing.T) {
	service, _ := setupFulcioTestService(t)

	ctx := context.Background()

	// Generate a key pair for testing
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	require.NoError(t, err)

	// Set up a fake CAClient for testing

	// Test that an error is returned for an invalid token
	_, err = getCert(ctx, key, service.client, "invalid_token")
	require.Error(t, err)

	// Test that an error is returned for a token without a subject
	token := generateTestToken("", "")
	_, err = getCert(ctx, key, service.client, token)
	require.Error(t, err)

	// Test that an error is returned if the key cannot be loaded
	key2, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	require.NoError(t, err)
	_, err = getCert(ctx, key2, service.client, token)
	require.Error(t, err)

	// Generate a token with email claim for testing
	token = generateTestToken("test@example.com", "")
	// Test that a certificate is returned for a valid token and key
	cert, err := getCert(ctx, key, service.client, token)
	require.NoError(t, err)
	require.NotNil(t, cert)

	// Generate a token with subject claim for testing
	token = generateTestToken("", "examplesubject")
	// Test that a certificate is returned for a valid token and key
	cert, err = getCert(ctx, key, service.client, token)
	require.NoError(t, err)
	require.NotNil(t, cert)
}

func TestSigner(t *testing.T) {
	// Setup dummy CA client service
	service, url := setupFulcioTestService(t)
	defer service.server.Stop()

	ctx := context.Background()

	// Create mock token
	token := generateTestToken("foo@example.com", "")

	//pasre url to get hostname
	hostname := strings.Split(url, ":")[0]
	port := strings.Split(url, ":")[1]

	// Call Signer to generate a signer
	provider := New(WithFulcioURL(fmt.Sprintf("http://%v:%v", hostname, port)), WithToken(token))
	signer, err := provider.Signer(ctx)
	require.NoError(t, err)

	// Ensure signer is not nil
	require.NotNil(t, signer)
	provider = New(WithFulcioURL("https://test"), WithToken(token))
	_, err = provider.Signer(ctx)
	//this should be a tranport err since we cant actually test on 443 which is the default
	require.ErrorContains(t, err, "lookup test")

	// Test signer with token read from file
	// NOTE: this function could be refactored to accept a fileSystem or io.Reader so reading the file can be mocked,
	// but unsure if this is the way we want to go for now
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	rootDir := filepath.Dir(filepath.Dir(wd))
	tp := filepath.Join(rootDir, "hack", "test.token")

	provider = New(WithFulcioURL(fmt.Sprintf("http://%v:%v", hostname, port)), WithTokenPath(tp))
	_, err = provider.Signer(ctx)
	require.NoError(t, err)

	// Test signer with both token read from file and raw token
	provider = New(WithFulcioURL(fmt.Sprintf("http://%v:%v", hostname, port)), WithTokenPath(tp), WithToken(token))
	_, err = provider.Signer(ctx)
	require.ErrorContains(t, err, "only one of --fulcio-token-path or --fulcio-raw-token can be used")
}

func generateCertChain(t *testing.T) []string {
	rootKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	rootCertTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Root CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	rootCertDER, err := x509.CreateCertificate(rand.Reader, &rootCertTemplate, &rootCertTemplate, &rootKey.PublicKey, rootKey)
	require.NoError(t, err)

	intermediateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	intermediateCertTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Intermediate",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	intermediateCertDER, err := x509.CreateCertificate(rand.Reader, &intermediateCertTemplate, &rootCertTemplate, &intermediateKey.PublicKey, rootKey)
	require.NoError(t, err)

	leafKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	leafCertTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Leaf",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365),
		KeyUsage:  x509.KeyUsageDigitalSignature,
	}
	leafCertDER, err := x509.CreateCertificate(rand.Reader, &leafCertTemplate, &intermediateCertTemplate, &leafKey.PublicKey, intermediateKey)
	require.NoError(t, err)

	certs := []string{
		string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootCertDER})),
		string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: intermediateCertDER})),
		string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: leafCertDER})),
	}

	return certs
}
