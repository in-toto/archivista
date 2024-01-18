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
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/registry"
	"github.com/in-toto/go-witness/signer"
)

const (
	defaultPkiSecretsEnginePath = "pki"
)

func init() {
	signer.Register("vault", func() signer.SignerProvider { return New() },
		registry.StringConfigOption(
			"url",
			"Base url of the Vault instance to connect to",
			"",
			func(sp signer.SignerProvider, url string) (signer.SignerProvider, error) {
				vsp, ok := sp.(*VaultSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a vault signer provider")
				}

				WithUrl(url)(vsp)
				return vsp, nil
			},
		),
		registry.StringConfigOption(
			"pki-secrets-engine-path",
			"Path to the Vault PKI Secrets Engine to use",
			defaultPkiSecretsEnginePath,
			func(sp signer.SignerProvider, pkiSecretsEnginePath string) (signer.SignerProvider, error) {
				vsp, ok := sp.(*VaultSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a vault signer provider")
				}

				WithPkiSecretsEnginePath(pkiSecretsEnginePath)(vsp)
				return vsp, nil
			},
		),

		registry.StringConfigOption(
			"token",
			"Token to use to connect to Vault",
			"",
			func(sp signer.SignerProvider, token string) (signer.SignerProvider, error) {
				vsp, ok := sp.(*VaultSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a vault signer provider")
				}

				WithToken(token)(vsp)
				return vsp, nil
			},
		),
		registry.StringConfigOption(
			"namespace",
			"Vault namespace to use",
			"",
			func(sp signer.SignerProvider, namespace string) (signer.SignerProvider, error) {
				vsp, ok := sp.(*VaultSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a vault signer provider")
				}

				WithNamespace(namespace)(vsp)
				return vsp, nil
			},
		),
		registry.StringConfigOption(
			"role",
			"Name of the Vault role to generate the certificate for",
			"",
			func(sp signer.SignerProvider, role string) (signer.SignerProvider, error) {
				vsp, ok := sp.(*VaultSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a vault signer provider")
				}

				WithRole(role)(vsp)
				return vsp, nil
			},
		),
		registry.StringConfigOption(
			"commonname",
			"Common name to use for the generated certificate. Must be allowed by the vault role policy",
			"",
			func(sp signer.SignerProvider, cn string) (signer.SignerProvider, error) {
				vsp, ok := sp.(*VaultSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a vault signer provider")
				}

				WithCommonName(cn)(vsp)
				return vsp, nil
			},
		),
		registry.StringSliceConfigOption(
			"altnames",
			"Alt names to use for the generated certificate. All alt names must be allowed by the vault role policy",
			[]string{},
			func(sp signer.SignerProvider, ans []string) (signer.SignerProvider, error) {
				vsp, ok := sp.(*VaultSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a vault signer provider")
				}

				WithAltNames(ans)(vsp)
				return vsp, nil
			},
		),
		registry.DurationConfigOption(
			"ttl",
			"Time to live for the generated certificate. Defaults to the vault role policy's configured TTL if not provided",
			time.Duration(0),
			func(sp signer.SignerProvider, ttl time.Duration) (signer.SignerProvider, error) {
				vsp, ok := sp.(*VaultSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a vault signer provider")
				}

				WithTtl(ttl)(vsp)
				return vsp, nil
			},
		),
	)
}

type VaultSignerProvider struct {
	requestIssuer        func(context.Context) (issueResponse, error)
	url                  string
	pkiSecretsEnginePath string
	token                string
	namespace            string
	role                 string
	commonName           string
	altNames             []string
	ttl                  time.Duration
}

type Option func(*VaultSignerProvider)

func WithUrl(url string) Option {
	return func(vsp *VaultSignerProvider) {
		vsp.url = url
	}
}

func WithPkiSecretsEnginePath(pkiSecretsEnginePath string) Option {
	return func(vsp *VaultSignerProvider) {
		vsp.pkiSecretsEnginePath = pkiSecretsEnginePath
	}
}

func WithToken(token string) Option {
	return func(vsp *VaultSignerProvider) {
		vsp.token = token
	}
}

func WithNamespace(namespace string) Option {
	return func(vsp *VaultSignerProvider) {
		vsp.namespace = namespace
	}
}

func WithRole(role string) Option {
	return func(vsp *VaultSignerProvider) {
		vsp.role = role
	}
}

func WithCommonName(cn string) Option {
	return func(vsp *VaultSignerProvider) {
		vsp.commonName = cn
	}
}

func WithAltNames(ans []string) Option {
	return func(vsp *VaultSignerProvider) {
		vsp.altNames = ans
	}
}

func WithTtl(ttl time.Duration) Option {
	return func(vsp *VaultSignerProvider) {
		vsp.ttl = ttl
	}
}

func New(opts ...Option) *VaultSignerProvider {
	vsp := VaultSignerProvider{}
	vsp.requestIssuer = vsp.requestCertificate

	for _, opt := range opts {
		opt(&vsp)
	}

	return &vsp
}

func (vsp *VaultSignerProvider) Signer(ctx context.Context) (cryptoutil.Signer, error) {
	if len(vsp.url) == 0 {
		return nil, fmt.Errorf("url is a required option")
	}

	if len(vsp.token) == 0 {
		return nil, fmt.Errorf("token is a required option")
	}

	if len(vsp.role) == 0 {
		return nil, fmt.Errorf("role is a required option")
	}

	resp, err := vsp.requestIssuer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to issue certificate: %w", err)
	}

	cert, err := cryptoutil.TryParseCertificate([]byte(resp.Data.Certificate))
	if err != nil {
		return nil, fmt.Errorf("could not parse certificate from response: %w", err)
	}

	intermediates := make([]*x509.Certificate, 0)
	for _, i := range resp.Data.CaChain {
		intermediate, err := cryptoutil.TryParseCertificate([]byte(i))
		if err != nil {
			return nil, fmt.Errorf("could not parse intermediate certificate from response: %w", err)
		}

		intermediates = append(intermediates, intermediate)
	}

	return cryptoutil.NewSignerFromReader(
		strings.NewReader(resp.Data.PrivateKey),
		cryptoutil.SignWithCertificate(cert),
		cryptoutil.SignWithIntermediates(intermediates),
	)
}
