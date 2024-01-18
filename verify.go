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

package witness

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"

	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/dsse"
	"github.com/in-toto/go-witness/policy"
	"github.com/in-toto/go-witness/source"
	"github.com/in-toto/go-witness/timestamp"
)

func VerifySignature(r io.Reader, verifiers ...cryptoutil.Verifier) (dsse.Envelope, error) {
	decoder := json.NewDecoder(r)
	envelope := dsse.Envelope{}
	if err := decoder.Decode(&envelope); err != nil {
		return envelope, fmt.Errorf("failed to parse dsse envelope: %w", err)
	}

	_, err := envelope.Verify(dsse.VerifyWithVerifiers(verifiers...))
	return envelope, err
}

type verifyOptions struct {
	policyTimestampAuthorities []dsse.TimestampVerifier
	policyCARoots              []*x509.Certificate
	policyCAIntermediates      []*x509.Certificate
	policyEnvelope             dsse.Envelope
	policyVerifiers            []cryptoutil.Verifier
	collectionSource           source.Sourcer
	subjectDigests             []string
}

type VerifyOption func(*verifyOptions)

func VerifyWithSubjectDigests(subjectDigests []cryptoutil.DigestSet) VerifyOption {
	return func(vo *verifyOptions) {
		for _, set := range subjectDigests {
			for _, digest := range set {
				vo.subjectDigests = append(vo.subjectDigests, digest)
			}
		}
	}
}

func VerifyWithCollectionSource(source source.Sourcer) VerifyOption {
	return func(vo *verifyOptions) {
		vo.collectionSource = source
	}
}

func VerifyWithPolicyTimestampAuthorities(authorities []dsse.TimestampVerifier) VerifyOption {
	return func(vo *verifyOptions) {
		vo.policyTimestampAuthorities = authorities
	}
}

func VerifyWithPolicyCARoots(roots []*x509.Certificate) VerifyOption {
	return func(vo *verifyOptions) {
		vo.policyCARoots = roots
	}
}

// Verify verifies a set of attestations against a provided policy. The set of attestations that satisfy the policy will be returned
// if verifiation is successful.
func Verify(ctx context.Context, policyEnvelope dsse.Envelope, policyVerifiers []cryptoutil.Verifier, opts ...VerifyOption) (map[string][]source.VerifiedCollection, error) {
	vo := verifyOptions{
		policyEnvelope:  policyEnvelope,
		policyVerifiers: policyVerifiers,
	}

	for _, opt := range opts {
		opt(&vo)
	}

	if _, err := vo.policyEnvelope.Verify(dsse.VerifyWithVerifiers(vo.policyVerifiers...), dsse.VerifyWithTimestampVerifiers(vo.policyTimestampAuthorities...), dsse.VerifyWithRoots(vo.policyCARoots...), dsse.VerifyWithIntermediates(vo.policyCAIntermediates...)); err != nil {
		return nil, fmt.Errorf("could not verify policy: %w", err)
	}

	pol := policy.Policy{}
	if err := json.Unmarshal(vo.policyEnvelope.Payload, &pol); err != nil {
		return nil, fmt.Errorf("failed to unmarshal policy from envelope: %w", err)
	}

	pubKeysById, err := pol.PublicKeyVerifiers()
	if err != nil {
		return nil, fmt.Errorf("failed to get pulic keys from policy: %w", err)
	}

	pubkeys := make([]cryptoutil.Verifier, 0)
	for _, pubkey := range pubKeysById {
		pubkeys = append(pubkeys, pubkey)
	}

	trustBundlesById, err := pol.TrustBundles()
	if err != nil {
		return nil, fmt.Errorf("failed to load policy trust bundles: %w", err)
	}

	roots := make([]*x509.Certificate, 0)
	intermediates := make([]*x509.Certificate, 0)
	for _, trustBundle := range trustBundlesById {
		roots = append(roots, trustBundle.Root)
		intermediates = append(intermediates, intermediates...)
	}

	timestampAuthoritiesById, err := pol.TimestampAuthorityTrustBundles()
	if err != nil {
		return nil, fmt.Errorf("failed to load policy timestamp authorities: %w", err)
	}

	timestampVerifiers := make([]dsse.TimestampVerifier, 0)
	for _, timestampAuthority := range timestampAuthoritiesById {
		certs := []*x509.Certificate{timestampAuthority.Root}
		certs = append(certs, timestampAuthority.Intermediates...)
		timestampVerifiers = append(timestampVerifiers, timestamp.NewVerifier(timestamp.VerifyWithCerts(certs)))
	}

	verifiedSource := source.NewVerifiedSource(
		vo.collectionSource,
		dsse.VerifyWithVerifiers(pubkeys...),
		dsse.VerifyWithRoots(roots...),
		dsse.VerifyWithIntermediates(intermediates...),
		dsse.VerifyWithTimestampVerifiers(timestampVerifiers...),
	)
	accepted, err := pol.Verify(ctx, policy.WithSubjectDigests(vo.subjectDigests), policy.WithVerifiedSource(verifiedSource))
	if err != nil {
		return nil, fmt.Errorf("failed to verify policy: %w", err)
	}

	return accepted, nil
}
