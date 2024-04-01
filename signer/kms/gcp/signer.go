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

package gcp

import (
	"context"
	"crypto"
	"fmt"
	"hash/crc32"
	"io"

	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/log"
	kms "github.com/in-toto/go-witness/signer/kms"
)

var gcpSupportedHashFuncs = []crypto.Hash{
	crypto.SHA256,
	crypto.SHA512,
	crypto.SHA384,
}

// SignerVerifier is a cryptoutil.SignerVerifier that uses the AWS Key Management Service
type SignerVerifier struct {
	reference string
	client    *gcpClient
	hashFunc  crypto.Hash
}

// LoadSignerVerifier generates signatures using the specified key object in AWS KMS and hash algorithm.
func LoadSignerVerifier(ctx context.Context, ksp *kms.KMSSignerProvider) (*SignerVerifier, error) {
	g := &SignerVerifier{
		reference: ksp.Reference,
	}

	var err error
	g.client, err = newGCPClient(ctx, ksp)
	if err != nil {
		return nil, err
	}

	for _, hashFunc := range gcpSupportedHashFuncs {
		if hashFunc == ksp.HashFunc {
			g.hashFunc = ksp.HashFunc
		}
	}

	if g.hashFunc == 0 {
		return nil, fmt.Errorf("unsupported hash function: %v", ksp.HashFunc)
	}

	return g, nil
}

// NOTE: This might be all wrong but setting it like so for now
//
// KeyID returns the key identifier for the key used by this signer.
func (g *SignerVerifier) KeyID() (string, error) {
	return g.reference, nil
}

// Sign signs the provided message using GCP KMS. If the message is provided,
// this method will compute the digest according to the hash function specified
// when the Signer was created.
func (g *SignerVerifier) Sign(message io.Reader) ([]byte, error) {
	var err error
	ctx := context.Background()
	var digest []byte

	var signerOpts crypto.SignerOpts
	signerOpts, err = g.client.getHashFunc()
	if err != nil {
		return nil, fmt.Errorf("failed to get default hash function: %w", err)
	}

	hf := signerOpts.HashFunc()

	digest, _, err = cryptoutil.ComputeDigest(message, hf, gcpSupportedHashFuncs)
	if err != nil {
		return nil, err
	}

	crc32cHasher := crc32.New(crc32.MakeTable(crc32.Castagnoli))
	_, err = crc32cHasher.Write(digest)
	if err != nil {
		return nil, err
	}

	return g.client.sign(ctx, digest, hf, crc32cHasher.Sum32())
}

// Verifier returns a cryptoutil.Verifier that can be used to verify signatures created by this signer.
func (g *SignerVerifier) Verifier() (cryptoutil.Verifier, error) {
	return g, nil
}

// PublicKey returns the public key that can be used to verify signatures created by
// this signer.
func (g *SignerVerifier) PublicKey(ctx context.Context) (crypto.PublicKey, error) {
	return g.client.public(ctx)
}

// Bytes returns the bytes of the public key that can be used to verify signatures created by the signer.
func (g *SignerVerifier) Bytes() ([]byte, error) {
	ckv, err := g.client.getCKV()
	if err != nil {
		return nil, fmt.Errorf("failed to get KMS key version: %w", err)
	}

	return cryptoutil.PublicPemBytes(ckv.PublicKey)
}

// VerifySignature verifies the signature for the given message, returning
// nil if the verification succeeded, and an error message otherwise.
func (g *SignerVerifier) Verify(message io.Reader, sig []byte) (err error) {
	err = g.client.verify(message, sig)
	if err != nil {
		log.Info(err.Error())
	}

	return err
}

// SupportedAlgorithms returns the list of algorithms supported by the AWS KMS service
func (*SignerVerifier) SupportedAlgorithms() (result []string) {
	for k := range algorithmMap {
		result = append(result, k)
	}
	return
}

// DefaultAlgorithm returns the default algorithm for the GCP KMS service
func (g *SignerVerifier) DefaultAlgorithm() string {
	return AlgorithmECDSAP256SHA256
}
