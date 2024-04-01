// Copyright 2024 The Witness Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/in-toto/go-witness/cryptoutil"
)

func CreateRsaKey() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	return priv, &priv.PublicKey, nil
}

func CreateTestKey() (cryptoutil.Signer, cryptoutil.Verifier, []byte, error) {
	privKey, _, err := CreateRsaKey()
	if err != nil {
		return nil, nil, nil, err
	}

	signer := cryptoutil.NewRSASigner(privKey, crypto.SHA256)
	verifier := cryptoutil.NewRSAVerifier(&privKey.PublicKey, crypto.SHA256)
	keyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, nil, nil, err
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: keyBytes})

	return signer, verifier, pemBytes, nil
}

func CreateCert(priv, pub interface{}, temp, parent *x509.Certificate) (*x509.Certificate, error) {
	var err error
	temp.SerialNumber, err = rand.Int(rand.Reader, big.NewInt(4294967295))
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, temp, parent, pub, priv)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(certBytes)
}

func CreateRoot() (*x509.Certificate, interface{}, error) {
	priv, pub, err := CreateRsaKey()
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		DNSNames: []string{"in-toto.io"},
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"in-toto"},
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

	cert, err := CreateCert(priv, pub, template, template)
	return cert, priv, err
}

func CreateIntermediate(parent *x509.Certificate, parentPriv interface{}) (*x509.Certificate, interface{}, error) {
	priv, pub, err := CreateRsaKey()
	if err != nil {
		return nil, nil, err
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

	cert, err := CreateCert(parentPriv, pub, template, parent)
	return cert, priv, err
}

func CreateLeaf(parent *x509.Certificate, parentPriv interface{}) (*x509.Certificate, interface{}, error) {
	priv, pub, err := CreateRsaKey()
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		DNSNames: []string{"in-toto.io"},
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"In-toto"},
			CommonName:   "Test Leaf",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	cert, err := CreateCert(parentPriv, pub, template, parent)
	return cert, priv, err
}
