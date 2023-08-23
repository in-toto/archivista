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

package source

import (
	"context"

	"github.com/testifysec/go-witness/cryptoutil"
	"github.com/testifysec/go-witness/dsse"
	"github.com/testifysec/go-witness/log"
)

type VerifiedCollection struct {
	Verifiers []cryptoutil.Verifier
	CollectionEnvelope
}

type VerifiedSourcer interface {
	Search(ctx context.Context, collectionName string, subjectDigests, attestations []string) ([]VerifiedCollection, error)
}

type VerifiedSource struct {
	source     Sourcer
	verifyOpts []dsse.VerificationOption
}

func NewVerifiedSource(source Sourcer, verifyOpts ...dsse.VerificationOption) *VerifiedSource {
	return &VerifiedSource{source, verifyOpts}
}

func (s *VerifiedSource) Search(ctx context.Context, collectionName string, subjectDigests, attestations []string) ([]VerifiedCollection, error) {
	unverified, err := s.source.Search(ctx, collectionName, subjectDigests, attestations)
	if err != nil {
		return nil, err
	}

	verified := make([]VerifiedCollection, 0)
	for _, toVerify := range unverified {
		envelopeVerifiers, err := toVerify.Envelope.Verify(s.verifyOpts...)
		if err != nil {
			log.Debugf("(verified source) skipping envelope: couldn't verify enveloper's signature with the policy's verifiers: %+v", err)
			continue
		}

		passedVerifiers := make([]cryptoutil.Verifier, 0)
		for _, verifier := range envelopeVerifiers {
			passedVerifiers = append(passedVerifiers, verifier.Verifier)
		}

		verified = append(verified, VerifiedCollection{
			Verifiers:          passedVerifiers,
			CollectionEnvelope: toVerify,
		})
	}

	return verified, nil
}
