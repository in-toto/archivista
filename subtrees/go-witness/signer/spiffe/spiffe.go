// Copyright 2021 The Witness Contributors
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

package spiffe

import (
	"context"
	"fmt"

	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"github.com/testifysec/go-witness/cryptoutil"
	"github.com/testifysec/go-witness/registry"
	"github.com/testifysec/go-witness/signer"
)

func init() {
	signer.Register("spiffe", func() signer.SignerProvider { return New() },
		registry.StringConfigOption(
			"socket-path",
			"Path to the SPIFFE Workload API Socket",
			"",
			func(sp signer.SignerProvider, socketPath string) (signer.SignerProvider, error) {
				ssp, ok := sp.(SpiffeSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a spiffe signer provider")
				}

				WithSocketPath(socketPath)(&ssp)
				return ssp, nil
			},
		),
	)
}

type SpiffeSignerProvider struct {
	SocketPath string
}

type ErrInvalidSVID string

func (e ErrInvalidSVID) Error() string {
	return fmt.Sprintf("invalid svid: %v", string(e))
}

type Option func(*SpiffeSignerProvider)

func WithSocketPath(socketPath string) Option {
	return func(ssp *SpiffeSignerProvider) {
		ssp.SocketPath = socketPath
	}
}

func New(opts ...Option) SpiffeSignerProvider {
	ssp := SpiffeSignerProvider{}
	for _, opt := range opts {
		opt(&ssp)
	}

	return ssp
}

func (ssp SpiffeSignerProvider) Signer(ctx context.Context) (cryptoutil.Signer, error) {
	if len(ssp.SocketPath) == 0 {
		return nil, fmt.Errorf("socker path cannot be empty")
	}

	svidCtx, err := workloadapi.FetchX509Context(ctx, workloadapi.WithAddr(ssp.SocketPath))
	if err != nil {
		return nil, err
	}

	svid := svidCtx.DefaultSVID()
	if len(svid.Certificates) <= 0 {
		return nil, ErrInvalidSVID("no certificates")
	}

	if svid.PrivateKey == nil {
		return nil, ErrInvalidSVID("no private key")
	}

	return cryptoutil.NewSigner(svid.PrivateKey, cryptoutil.SignWithIntermediates(svid.Certificates[1:]), cryptoutil.SignWithCertificate(svid.Certificates[0]))
}
