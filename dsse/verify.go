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

package dsse

import (
	"bytes"
	"context"
	"crypto/x509"
	"io"
	"time"

	"github.com/testifysec/go-witness/cryptoutil"
)

type TimestampVerifier interface {
	Verify(context.Context, io.Reader, io.Reader) (time.Time, error)
}

type verificationOptions struct {
	roots              []*x509.Certificate
	intermediates      []*x509.Certificate
	verifiers          []cryptoutil.Verifier
	threshold          int
	timestampVerifiers []TimestampVerifier
}

type VerificationOption func(*verificationOptions)

func VerifyWithRoots(roots ...*x509.Certificate) VerificationOption {
	return func(vo *verificationOptions) {
		vo.roots = roots
	}
}

func VerifyWithIntermediates(intermediates ...*x509.Certificate) VerificationOption {
	return func(vo *verificationOptions) {
		vo.intermediates = intermediates
	}
}

func VerifyWithVerifiers(verifiers ...cryptoutil.Verifier) VerificationOption {
	return func(vo *verificationOptions) {
		vo.verifiers = verifiers
	}
}

func VerifyWithThreshold(threshold int) VerificationOption {
	return func(vo *verificationOptions) {
		vo.threshold = threshold
	}
}

func VerifyWithTimestampVerifiers(verifiers ...TimestampVerifier) VerificationOption {
	return func(vo *verificationOptions) {
		vo.timestampVerifiers = verifiers
	}
}

type PassedVerifier struct {
	Verifier                 cryptoutil.Verifier
	PassedTimestampVerifiers []TimestampVerifier
}

func (e Envelope) Verify(opts ...VerificationOption) ([]PassedVerifier, error) {
	options := &verificationOptions{
		threshold: 1,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.threshold <= 0 {
		return nil, ErrInvalidThreshold(options.threshold)
	}

	pae := preauthEncode(e.PayloadType, e.Payload)
	if len(e.Signatures) == 0 {
		return nil, ErrNoSignatures{}
	}

	matchingSigFound := false
	passedVerifiers := make([]PassedVerifier, 0)
	for _, sig := range e.Signatures {
		if sig.Certificate != nil && len(sig.Certificate) > 0 {
			cert, err := cryptoutil.TryParseCertificate(sig.Certificate)
			if err != nil {
				continue
			}

			sigIntermediates := make([]*x509.Certificate, 0)
			for _, int := range sig.Intermediates {
				intCert, err := cryptoutil.TryParseCertificate(int)
				if err != nil {
					continue
				}

				sigIntermediates = append(sigIntermediates, intCert)
			}

			sigIntermediates = append(sigIntermediates, options.intermediates...)
			if len(options.timestampVerifiers) == 0 {
				if verifier, err := verifyX509Time(cert, sigIntermediates, options.roots, pae, sig.Signature, time.Now()); err == nil {
					matchingSigFound = true
					passedVerifiers = append(passedVerifiers, PassedVerifier{Verifier: verifier})
				}
			} else {
				var passedVerifier cryptoutil.Verifier
				passedTimestampVerifiers := []TimestampVerifier{}

				for _, timestampVerifier := range options.timestampVerifiers {
					for _, sigTimestamp := range sig.Timestamps {
						timestamp, err := timestampVerifier.Verify(context.TODO(), bytes.NewReader(sigTimestamp.Data), bytes.NewReader(sig.Signature))
						if err != nil {
							continue
						}

						if verifier, err := verifyX509Time(cert, sigIntermediates, options.roots, pae, sig.Signature, timestamp); err == nil {
							passedVerifier = verifier
							passedTimestampVerifiers = append(passedTimestampVerifiers, timestampVerifier)
						}
					}
				}

				if len(passedTimestampVerifiers) > 0 {
					matchingSigFound = true
					passedVerifiers = append(passedVerifiers, PassedVerifier{
						Verifier:                 passedVerifier,
						PassedTimestampVerifiers: passedTimestampVerifiers,
					})
				}
			}
		}

		for _, verifier := range options.verifiers {
			if verifier != nil {
				if err := verifier.Verify(bytes.NewReader(pae), sig.Signature); err == nil {
					passedVerifiers = append(passedVerifiers, PassedVerifier{Verifier: verifier})
					matchingSigFound = true
				}
			}
		}
	}

	if !matchingSigFound {
		return nil, ErrNoMatchingSigs{}
	}

	if len(passedVerifiers) < options.threshold {
		return passedVerifiers, ErrThresholdNotMet{Theshold: options.threshold, Acutal: len(passedVerifiers)}
	}

	return passedVerifiers, nil
}

func verifyX509Time(cert *x509.Certificate, sigIntermediates, roots []*x509.Certificate, pae, sig []byte, trustedTime time.Time) (cryptoutil.Verifier, error) {
	verifier, err := cryptoutil.NewX509Verifier(cert, sigIntermediates, roots, trustedTime)
	if err != nil {
		return nil, err
	}

	if err := verifier.Verify(bytes.NewReader(pae), sig); err != nil {
		return nil, err
	}

	return verifier, nil
}
