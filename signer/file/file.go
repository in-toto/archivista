// Copyright 2022 The Witness Contributors
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

package file

import (
	"context"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/registry"
	"github.com/in-toto/go-witness/signer"
)

func init() {
	signer.Register("file", func() signer.SignerProvider { return New() },
		registry.StringConfigOption(
			"key-path",
			"Path to the file containing the private key",
			"",
			func(sp signer.SignerProvider, keyPath string) (signer.SignerProvider, error) {
				ksp, ok := sp.(FileSignerProvider)
				if !ok {
					return ksp, fmt.Errorf("provided signer provider is not a file signer provider")
				}

				WithKeyPath(keyPath)(&ksp)
				return ksp, nil
			},
		),
		registry.StringConfigOption(
			"cert-path",
			"Path to the file containing the certificate for the private key",
			"",
			func(sp signer.SignerProvider, certPath string) (signer.SignerProvider, error) {
				ksp, ok := sp.(FileSignerProvider)
				if !ok {
					return ksp, fmt.Errorf("provided signer provider is not a file signer provider")
				}

				WithCertPath(certPath)(&ksp)
				return ksp, nil
			},
		),
		registry.StringSliceConfigOption(
			"intermediate-paths",
			"Paths to files containing intermediates required to establish trust of the signer's certificate to a root",
			[]string{},
			func(sp signer.SignerProvider, intermediatePaths []string) (signer.SignerProvider, error) {
				ksp, ok := sp.(FileSignerProvider)
				if !ok {
					return ksp, fmt.Errorf("provided signer provider is not a file signer provider")
				}

				WithIntermediatePaths(intermediatePaths)(&ksp)
				return ksp, nil
			},
		),
	)
}

type FileSignerProvider struct {
	KeyPath           string
	CertPath          string
	IntermediatePaths []string
}

type Option func(fsp *FileSignerProvider)

func WithKeyPath(keyPath string) Option {
	return func(fsp *FileSignerProvider) {
		fsp.KeyPath = keyPath
	}
}

func WithCertPath(certPath string) Option {
	return func(fsp *FileSignerProvider) {
		fsp.CertPath = certPath
	}
}

func WithIntermediatePaths(intermediatePaths []string) Option {
	return func(fsp *FileSignerProvider) {
		fsp.IntermediatePaths = intermediatePaths
	}
}

func New(opts ...Option) FileSignerProvider {
	fsp := FileSignerProvider{}
	for _, opt := range opts {
		opt(&fsp)
	}

	return fsp
}

func (fsp FileSignerProvider) Signer(ctx context.Context) (cryptoutil.Signer, error) {
	keyFile, err := os.Open(fsp.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open key file: %w", err)
	}

	defer keyFile.Close()
	key, err := cryptoutil.TryParseKeyFromReader(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load key: %w", err)
	}

	signerOpts := []cryptoutil.SignerOption{}
	if fsp.CertPath != "" {
		leaf, err := loadCert(fsp.CertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load certificate: %w", err)
		}

		signerOpts = append(signerOpts, cryptoutil.SignWithCertificate(leaf))
	}

	if len(fsp.IntermediatePaths) > 0 {
		intermediates := []*x509.Certificate{}
		for _, path := range fsp.IntermediatePaths {
			cert, err := loadCert(path)
			if err != nil {
				return nil, fmt.Errorf("failed to load intermediate: %w", err)
			}

			intermediates = append(intermediates, cert)
		}

		signerOpts = append(signerOpts, cryptoutil.SignWithIntermediates(intermediates))
	}

	return cryptoutil.NewSigner(key, signerOpts...)
}

func loadCert(path string) (*x509.Certificate, error) {
	certFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	defer certFile.Close()
	possibleCert, err := cryptoutil.TryParseKeyFromReader(certFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate")
	}

	cert, ok := possibleCert.(*x509.Certificate)
	if !ok {
		return nil, fmt.Errorf("%v is not a x509 certificate", path)
	}

	return cert, nil
}
