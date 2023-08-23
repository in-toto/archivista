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

package vault

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testifysec/go-witness/cryptoutil"
)

func createRsaKey() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	return priv, &priv.PublicKey, nil
}

func createCert(priv, pub interface{}, temp, parent *x509.Certificate) (*x509.Certificate, []byte, error) {
	var err error
	temp.SerialNumber, err = rand.Int(rand.Reader, big.NewInt(4294967295))
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, temp, parent, pub, priv)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(certBytes)
	return cert, certBytes, err
}

func createRoot() (*x509.Certificate, []byte, *rsa.PrivateKey, error) {
	priv, pub, err := createRsaKey()
	if err != nil {
		return nil, nil, nil, err
	}

	template := &x509.Certificate{
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"TestifySec"},
			CommonName:   "Test Root",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        false,
		MaxPathLen:            2,
	}

	cert, certBytes, err := createCert(priv, pub, template, template)
	return cert, certBytes, priv, err
}

func createIntermediate(parent *x509.Certificate, parentPriv interface{}) (*x509.Certificate, []byte, *rsa.PrivateKey, error) {
	priv, pub, err := createRsaKey()
	if err != nil {
		return nil, nil, nil, err
	}

	template := &x509.Certificate{
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"TestifySec"},
			CommonName:   "Test Intermediate",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        false,
		MaxPathLen:            1,
	}

	cert, certBytes, err := createCert(parentPriv, pub, template, parent)
	return cert, certBytes, priv, err
}

func createLeaf(parent *x509.Certificate, parentPriv interface{}) (*x509.Certificate, []byte, *rsa.PrivateKey, error) {
	priv, pub, err := createRsaKey()
	if err != nil {
		return nil, nil, nil, err
	}

	template := &x509.Certificate{
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"TestifySec"},
			CommonName:   "Test Leaf",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	cert, certBytes, err := createCert(parentPriv, pub, template, parent)
	return cert, certBytes, priv, err
}

func pemString(t string, b []byte) string {
	pemBlock := &pem.Block{
		Type:  t,
		Bytes: b,
	}

	pemBytes := pem.EncodeToMemory(pemBlock)
	return string(pemBytes)
}

func TestSigner(t *testing.T) {
	sp := New(WithToken("test token"), WithUrl("test url"), WithRole("dummy role"))
	root, _, rootPriv, err := createRoot()
	require.NoError(t, err)
	intermediate, intermediateBytes, intPriv, err := createIntermediate(root, rootPriv)
	require.NoError(t, err)
	_, leafBytes, priv, err := createLeaf(intermediate, intPriv)
	require.NoError(t, err)

	sp.requestIssuer = func(ctx context.Context) (issueResponse, error) {
		return issueResponse{
			Data: issueResponseData{
				Certificate: pemString("CERTIFICATE", leafBytes),
				CaChain: []string{
					pemString("CERTIFICATE", intermediateBytes),
				},
				PrivateKey: pemString("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(priv)),
			},
		}, nil
	}

	signer, err := sp.Signer(context.Background())
	require.NoError(t, err)
	assert.IsType(t, signer, &cryptoutil.X509Signer{})
	x509Signer := signer.(*cryptoutil.X509Signer)
	assert.Len(t, x509Signer.Intermediates(), 1)
}

func TestValidation(t *testing.T) {
	t.Run("without url", func(t *testing.T) {
		sp := New()
		_, err := sp.Signer(context.Background())
		assert.ErrorContains(t, err, "url is a required option")
	})
	t.Run("without token", func(t *testing.T) {
		sp := New(WithUrl("test"))
		_, err := sp.Signer(context.Background())
		assert.ErrorContains(t, err, "token is a required option")
	})
	t.Run("without role", func(t *testing.T) {
		sp := New(WithUrl("test"), WithToken("test"))
		_, err := sp.Signer(context.Background())
		assert.ErrorContains(t, err, "role is a required option")
	})

}
