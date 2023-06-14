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

package policy

import (
	"bytes"
	"context"
	"crypto/x509"
	"time"

	"github.com/testifysec/go-witness/attestation"
	"github.com/testifysec/go-witness/cryptoutil"
	"github.com/testifysec/go-witness/log"
	"github.com/testifysec/go-witness/source"
)

const PolicyPredicate = "https://witness.testifysec.com/policy/v0.1"

type Policy struct {
	Expires              time.Time            `json:"expires"`
	Roots                map[string]Root      `json:"roots,omitempty"`
	TimestampAuthorities map[string]Root      `json:"timestampauthorities,omitempty"`
	PublicKeys           map[string]PublicKey `json:"publickeys,omitempty"`
	Steps                map[string]Step      `json:"steps"`
}

type Root struct {
	Certificate   []byte   `json:"certificate"`
	Intermediates [][]byte `json:"intermediates,omitempty"`
}

type PublicKey struct {
	KeyID string `json:"keyid"`
	Key   []byte `json:"key"`
}

// PublicKeyVerifiers returns verifiers for each of the policy's embedded public keys grouped by the key's ID
func (p Policy) PublicKeyVerifiers() (map[string]cryptoutil.Verifier, error) {
	verifiers := make(map[string]cryptoutil.Verifier)
	for _, key := range p.PublicKeys {
		verifier, err := cryptoutil.NewVerifierFromReader(bytes.NewReader(key.Key))
		if err != nil {
			return nil, err
		}

		keyID, err := verifier.KeyID()
		if err != nil {
			return nil, err
		}

		if keyID != key.KeyID {
			return nil, ErrKeyIDMismatch{
				Expected: key.KeyID,
				Actual:   keyID,
			}
		}

		verifiers[keyID] = verifier
	}

	return verifiers, nil
}

type TrustBundle struct {
	Root          *x509.Certificate
	Intermediates []*x509.Certificate
}

// TrustBundles returns the policy's x509 roots and intermediates grouped by the root's ID
func (p Policy) TrustBundles() (map[string]TrustBundle, error) {
	return trustBundlesFromRoots(p.Roots)
}

func (p Policy) TimestampAuthorityTrustBundles() (map[string]TrustBundle, error) {
	return trustBundlesFromRoots(p.TimestampAuthorities)
}

func trustBundlesFromRoots(roots map[string]Root) (map[string]TrustBundle, error) {
	bundles := make(map[string]TrustBundle)
	for id, root := range roots {
		bundle := TrustBundle{}
		var err error
		bundle.Root, err = cryptoutil.TryParseCertificate(root.Certificate)
		if err != nil {
			return bundles, err
		}

		for _, intBytes := range root.Intermediates {
			cert, err := cryptoutil.TryParseCertificate(intBytes)
			if err != nil {
				return bundles, err
			}

			bundle.Intermediates = append(bundle.Intermediates, cert)
		}

		bundles[id] = bundle
	}

	return bundles, nil
}

type VerifyOption func(*verifyOptions)

type verifyOptions struct {
	verifiedSource source.VerifiedSourcer
	subjectDigests []string
	searchDepth    int
}

func WithVerifiedSource(verifiedSource source.VerifiedSourcer) VerifyOption {
	return func(vo *verifyOptions) {
		vo.verifiedSource = verifiedSource
	}
}

func WithSubjectDigests(subjectDigests []string) VerifyOption {
	return func(vo *verifyOptions) {
		vo.subjectDigests = subjectDigests
	}
}

func WithSearchDepth(depth int) VerifyOption {
	return func(vo *verifyOptions) {
		vo.searchDepth = depth
	}
}

func checkVerifyOpts(vo *verifyOptions) error {
	if vo.verifiedSource == nil {
		return ErrInvalidOption{
			Option: "verified source",
			Reason: "a verified attestation source is required",
		}
	}

	if len(vo.subjectDigests) == 0 {
		return ErrInvalidOption{
			Option: "subject digests",
			Reason: "at least one subject digest is required",
		}
	}

	if vo.searchDepth < 1 {
		return ErrInvalidOption{
			Option: "search depth",
			Reason: "search depth must be at least 1",
		}
	}

	return nil
}

func (p Policy) Verify(ctx context.Context, opts ...VerifyOption) (map[string][]source.VerifiedCollection, error) {
	vo := &verifyOptions{
		searchDepth: 3,
	}

	for _, opt := range opts {
		opt(vo)
	}

	if err := checkVerifyOpts(vo); err != nil {
		return nil, err
	}

	if time.Now().After(p.Expires) {
		return nil, ErrPolicyExpired(p.Expires)
	}

	trustBundles, err := p.TrustBundles()
	if err != nil {
		return nil, err
	}

	attestationsByStep := make(map[string][]string)
	for name, step := range p.Steps {
		for _, attestation := range step.Attestations {
			attestationsByStep[name] = append(attestationsByStep[name], attestation.Type)
		}
	}

	passedByStep := make(map[string][]source.VerifiedCollection)
	for depth := 0; depth < vo.searchDepth; depth++ {
		for stepName, step := range p.Steps {
			statements, err := vo.verifiedSource.Search(ctx, stepName, vo.subjectDigests, attestationsByStep[stepName])
			if err != nil {
				return nil, err
			}

			approvedCollections := step.checkFunctionaries(statements, trustBundles)
			stepResults := step.validateAttestations(approvedCollections)
			passedByStep[stepName] = append(passedByStep[stepName], stepResults.Passed...)
			for _, coll := range stepResults.Passed {
				for _, digestSet := range coll.Collection.BackRefs() {
					for _, digest := range digestSet {
						vo.subjectDigests = append(vo.subjectDigests, digest)
					}
				}
			}
		}

		if accepted, err := p.verifyArtifacts(passedByStep); err == nil {
			return accepted, nil
		}
	}

	return nil, ErrPolicyDenied{Reasons: []string{"failed to find set of attestations that satisfies the policy"}}
}

// checkFunctionaries checks to make sure the signature on each statement corresponds to a trusted functionary for
// the step the statement corresponds to
func (step Step) checkFunctionaries(verifiedStatements []source.VerifiedCollection, trustBundles map[string]TrustBundle) []source.VerifiedCollection {
	collections := make([]source.VerifiedCollection, 0)
	for _, verifiedStatement := range verifiedStatements {
		if verifiedStatement.Statement.PredicateType != attestation.CollectionType {
			log.Debugf("(policy) skipping statement: predicate type is not a collection (%v)", verifiedStatement.Statement.PredicateType)
			continue
		}

		for _, verifier := range verifiedStatement.Verifiers {
			verifierID, err := verifier.KeyID()
			if err != nil {
				log.Debugf("(policy) skipping verifier: could not get key id: %v", err)
				continue
			}

			for _, functionary := range step.Functionaries {
				if functionary.PublicKeyID != "" && functionary.PublicKeyID == verifierID {
					collections = append(collections, verifiedStatement)
					break
				}

				x509Verifier, ok := verifier.(*cryptoutil.X509Verifier)
				if !ok {
					log.Debugf("(policy) skipping verifier: verifier with ID %v is not a public key verifier or a x509 verifier", verifierID)
					continue
				}

				if len(functionary.CertConstraint.Roots) == 0 {
					log.Debugf("(policy) skipping verifier: verifier with ID %v is an x509 verifier, but step %v does not have any truested roots", verifierID, step)
					continue
				}

				if err := functionary.CertConstraint.Check(x509Verifier, trustBundles); err != nil {
					log.Debugf("(policy) skipping verifier: verifier with ID %v doesn't meet certificate constraint: %w", verifierID, err)
					continue
				}

				collections = append(collections, verifiedStatement)
			}
		}
	}

	return collections
}

// verifyArtifacts will check the artifacts (materials+products) of the step referred to by `ArtifactsFrom` against the
// materials of the original step.  This ensures file integrity between each step.
func (p Policy) verifyArtifacts(collectionsByStep map[string][]source.VerifiedCollection) (map[string][]source.VerifiedCollection, error) {
	acceptedByStep := make(map[string][]source.VerifiedCollection)
	for _, step := range p.Steps {
		accepted := make([]source.VerifiedCollection, 0)
		for _, collection := range collectionsByStep[step.Name] {
			if err := verifyCollectionArtifacts(step, collection, collectionsByStep); err == nil {
				accepted = append(accepted, collection)
			}
		}

		acceptedByStep[step.Name] = accepted
		if len(accepted) <= 0 {
			return nil, ErrNoAttestations(step.Name)
		}
	}

	return acceptedByStep, nil
}

func verifyCollectionArtifacts(step Step, collection source.VerifiedCollection, collectionsByStep map[string][]source.VerifiedCollection) error {
	mats := collection.Collection.Materials()
	for _, artifactsFrom := range step.ArtifactsFrom {
		accepted := make([]source.VerifiedCollection, 0)
		for _, testCollection := range collectionsByStep[artifactsFrom] {
			if err := compareArtifacts(mats, testCollection.Collection.Artifacts()); err != nil {
				break
			}

			accepted = append(accepted, testCollection)
		}

		if len(accepted) <= 0 {
			return ErrNoAttestations(step.Name)
		}
	}

	return nil
}

func compareArtifacts(mats map[string]cryptoutil.DigestSet, arts map[string]cryptoutil.DigestSet) error {
	for path, mat := range mats {
		art, ok := arts[path]
		if !ok {
			continue
		}

		if !mat.Equal(art) {
			return ErrMismatchArtifact{
				Artifact: art,
				Material: mat,
				Path:     path,
			}
		}
	}

	return nil
}
