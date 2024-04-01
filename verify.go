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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/dsse"
	"github.com/in-toto/go-witness/log"
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
	policyTimestampAuthorities []timestamp.TimestampVerifier
	policyCARoots              []*x509.Certificate
	policyCAIntermediates      []*x509.Certificate
	policyCommonName           string
	policyDNSNames             []string
	policyEmails               []string
	policyOrganizations        []string
	policyURIs                 []string
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

func VerifyWithPolicyTimestampAuthorities(authorities []timestamp.TimestampVerifier) VerifyOption {
	return func(vo *verifyOptions) {
		vo.policyTimestampAuthorities = authorities
	}
}

func VerifyWithPolicyCARoots(roots []*x509.Certificate) VerifyOption {
	return func(vo *verifyOptions) {
		vo.policyCARoots = roots
	}
}

func VerifyWithPolicyCAIntermediates(intermediates []*x509.Certificate) VerifyOption {
	return func(vo *verifyOptions) {
		vo.policyCAIntermediates = intermediates
	}
}

func VerifyWithPolicyCertConstraints(commonName string, dnsNames []string, emails []string, organizations []string, uris []string) VerifyOption {
	return func(vo *verifyOptions) {
		vo.policyCommonName = commonName
		vo.policyDNSNames = dnsNames
		vo.policyEmails = emails
		vo.policyOrganizations = organizations
		vo.policyURIs = uris
	}
}

// Verify verifies a set of attestations against a provided policy. The set of attestations that satisfy the policy will be returned
// if verifiation is successful.
func Verify(ctx context.Context, policyEnvelope dsse.Envelope, policyVerifiers []cryptoutil.Verifier, opts ...VerifyOption) (map[string][]source.VerifiedCollection, error) {
	vo := verifyOptions{
		policyEnvelope:      policyEnvelope,
		policyVerifiers:     policyVerifiers,
		policyCommonName:    "*",
		policyDNSNames:      []string{"*"},
		policyOrganizations: []string{"*"},
		policyURIs:          []string{"*"},
		policyEmails:        []string{"*"},
	}

	for _, opt := range opts {
		opt(&vo)
	}

	if err := verifyPolicySignature(ctx, vo); err != nil {
		return nil, fmt.Errorf("failed to verify policy signature: %w", err)
	}

	log.Info("policy signature verified")

	pol := policy.Policy{}
	if err := json.Unmarshal(vo.policyEnvelope.Payload, &pol); err != nil {
		return nil, fmt.Errorf("failed to unmarshal policy from envelope: %w", err)
	}

	pubKeysById, err := pol.PublicKeyVerifiers()
	if err != nil {
		return nil, fmt.Errorf("failed to get public keys from policy: %w", err)
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

	timestampVerifiers := make([]timestamp.TimestampVerifier, 0)
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

func verifyPolicySignature(ctx context.Context, vo verifyOptions) error {
	passedPolicyVerifiers, err := vo.policyEnvelope.Verify(dsse.VerifyWithVerifiers(vo.policyVerifiers...), dsse.VerifyWithTimestampVerifiers(vo.policyTimestampAuthorities...), dsse.VerifyWithRoots(vo.policyCARoots...), dsse.VerifyWithIntermediates(vo.policyCAIntermediates...))
	if err != nil {
		return fmt.Errorf("could not verify policy: %w", err)
	}

	var passed bool
	for _, verifier := range passedPolicyVerifiers {
		kid, err := verifier.Verifier.KeyID()
		if err != nil {
			return fmt.Errorf("could not get verifier key id: %w", err)
		}

		var f policy.Functionary
		trustBundle := make(map[string]policy.TrustBundle)
		if _, ok := verifier.Verifier.(*cryptoutil.X509Verifier); ok {
			rootIDs := make([]string, 0)
			for _, root := range vo.policyCARoots {
				id := base64.StdEncoding.EncodeToString(root.Raw)
				rootIDs = append(rootIDs, id)
				trustBundle[id] = policy.TrustBundle{
					Root: root,
				}
			}

			f = policy.Functionary{
				Type: "root",
				CertConstraint: policy.CertConstraint{
					Roots:         rootIDs,
					CommonName:    vo.policyCommonName,
					URIs:          vo.policyURIs,
					Emails:        vo.policyEmails,
					Organizations: vo.policyOrganizations,
					DNSNames:      vo.policyDNSNames,
				},
			}

		} else {
			f = policy.Functionary{
				Type:        "key",
				PublicKeyID: kid,
			}
		}

		err = f.Validate(verifier.Verifier, trustBundle)
		if err != nil {
			log.Debugf("Policy Verifier %s failed failed to match supplied constraints: %w, continuing...", kid, err)
			continue
		}
		passed = true
	}

	if !passed {
		return fmt.Errorf("no policy verifiers passed verification")
	} else {
		return nil
	}
}
