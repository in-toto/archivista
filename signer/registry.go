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

package signer

import (
	"context"

	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/registry"
)

var (
	signerRegistry   = registry.New[SignerProvider]()
	verifierRegistry = registry.New[VerifierProvider]()
)

type SignerProvider interface {
	Signer(context.Context) (cryptoutil.Signer, error)
}

func Register(name string, factory func() SignerProvider, opts ...registry.Configurer) {
	signerRegistry.Register(name, factory, opts...)
}

func RegistryEntries() []registry.Entry[SignerProvider] {
	return signerRegistry.AllEntries()
}

func NewSignerProvider(name string, opts ...func(SignerProvider) (SignerProvider, error)) (SignerProvider, error) {
	return signerRegistry.NewEntity(name, opts...)
}

// NOTE: This is a temporary interface, and should not be used. It will be deprecated in a future release.
// The same applies to the functions that use this interface.
type VerifierProvider interface {
	Verifier(context.Context) (cryptoutil.Verifier, error)
}

func RegisterVerifier(name string, factory func() VerifierProvider, opts ...registry.Configurer) {
	verifierRegistry.Register(name, factory, opts...)
}

func VerifierRegistryEntries() []registry.Entry[VerifierProvider] {
	return verifierRegistry.AllEntries()
}

func NewVerifierProvider(name string, opts ...func(VerifierProvider) (VerifierProvider, error)) (VerifierProvider, error) {
	return verifierRegistry.NewEntity(name, opts...)
}
