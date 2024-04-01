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

package aws

import (
	"bytes"
	"context"
	"crypto"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/in-toto/go-witness/cryptoutil"
	kms "github.com/in-toto/go-witness/signer/kms"
)

var awsSupportedAlgorithms = []types.CustomerMasterKeySpec{
	types.CustomerMasterKeySpecRsa2048,
	types.CustomerMasterKeySpecRsa3072,
	types.CustomerMasterKeySpecRsa4096,
	types.CustomerMasterKeySpecEccNistP256,
	types.CustomerMasterKeySpecEccNistP384,
	types.CustomerMasterKeySpecEccNistP521,
}

var awsSupportedHashFuncs = []crypto.Hash{
	crypto.SHA256,
	crypto.SHA384,
	crypto.SHA512,
}

// SignerVerifier is a cryptoutil.SignerVerifier that uses the AWS Key Management Service
type SignerVerifier struct {
	reference string
	client    client
	hashFunc  crypto.Hash
}

// LoadSignerVerifier generates signatures using the specified key object in AWS KMS and hash algorithm.
func LoadSignerVerifier(ctx context.Context, ksp *kms.KMSSignerProvider) (*SignerVerifier, error) {
	a := &SignerVerifier{
		reference: ksp.Reference,
	}

	var err error
	a.client, err = newAWSClient(ctx, ksp)
	if err != nil {
		return nil, err
	}

	for _, hashFunc := range awsSupportedHashFuncs {
		if hashFunc == ksp.HashFunc {
			a.hashFunc = ksp.HashFunc
		}
	}

	if a.hashFunc == 0 {
		return nil, fmt.Errorf("unsupported hash function: %v", ksp.HashFunc)
	}

	return a, nil
}

// NOTE: This might be all wrong but setting it like so for now
//
// KeyID returns the key identifier for the key used by this signer.
func (a *SignerVerifier) KeyID() (string, error) {
	return a.reference, nil
}

// Sign signs the provided message using AWS KMS. If the message is provided,
// this method will compute the digest according to the hash function specified
// when the Signer was created.
func (a *SignerVerifier) Sign(message io.Reader) ([]byte, error) {
	var err error
	ctx := context.TODO()
	var digest []byte

	var signerOpts crypto.SignerOpts
	signerOpts, err = a.client.getHashFunc(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting fetching default hash function: %w", err)
	}

	hf := signerOpts.HashFunc()

	digest, _, err = cryptoutil.ComputeDigest(message, hf, awsSupportedHashFuncs)
	if err != nil {
		return nil, err
	}

	return a.client.sign(ctx, digest, hf)
}

// Verifier returns a cryptoutil.Verifier that can be used to verify signatures created by this signer.
func (a *SignerVerifier) Verifier() (cryptoutil.Verifier, error) {
	return a, nil
}

// Bytes returns the bytes of the public key that can be used to verify signatures created by the signer.
func (a *SignerVerifier) Bytes() ([]byte, error) {
	ctx := context.TODO()
	p, err := a.client.fetchPublicKey(ctx)
	if err != nil {
		return nil, err
	}

	return cryptoutil.PublicPemBytes(p)
}

// Verify verifies the signature for the given message, returning
// nil if the verification succeeded, and an error message otherwise.
func (a *SignerVerifier) Verify(message io.Reader, sig []byte) (err error) {
	ctx := context.TODO()

	return a.client.verify(ctx, bytes.NewReader(sig), message)
}

// SupportedAlgorithms returns the list of algorithms supported by the AWS KMS service
func (*SignerVerifier) SupportedAlgorithms() []string {
	s := make([]string, len(awsSupportedAlgorithms))
	for i := range awsSupportedAlgorithms {
		s[i] = string(awsSupportedAlgorithms[i])
	}
	return s
}

// DefaultAlgorithm returns the default algorithm for the AWS KMS service
func (*SignerVerifier) DefaultAlgorithm() string {
	return string(types.CustomerMasterKeySpecEccNistP256)
}
