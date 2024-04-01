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
	"crypto/rand"
	"crypto/rsa"
	"time"

	akms "github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/signer/kms"
	ttlcache "github.com/jellydator/ttlcache/v3"
)

var (
	aid = "012345678901"
	arn = "arn:aws:kms:us-west-2:012345678901:key/12345678-1234-1234-1234-123456789012"
)

func createRsaKey() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	return privKey, &privKey.PublicKey, nil
}

func createTestKey() (cryptoutil.Signer, cryptoutil.Verifier, error) {
	privKey, pubKey, err := createRsaKey()
	if err != nil {
		return nil, nil, err
	}

	signer := cryptoutil.NewRSASigner(privKey, crypto.SHA256)
	verifier := cryptoutil.NewRSAVerifier(pubKey, crypto.SHA256)
	return signer, verifier, nil
}

type fakeAWSClient struct {
	client     *akms.Client
	endpoint   string
	keyID      string
	alias      string
	keyCache   *ttlcache.Cache[string, cmk]
	privateKey *rsa.PrivateKey
	hash       crypto.Hash
}

func newFakeAWSClient(ctx context.Context, ksp *kms.KMSSignerProvider) (*fakeAWSClient, error) {
	a, err := newAWSClient(ctx, ksp)
	if err != nil {
		return nil, err
	}

	c := &fakeAWSClient{
		client:   a.client,
		endpoint: a.endpoint,
		keyID:    a.keyID,
		alias:    a.alias,
		keyCache: a.keyCache,
		hash:     ksp.HashFunc,
	}

	return c, nil
}

func (a *fakeAWSClient) fetchCMK(ctx context.Context) (*cmk, error) {
	var err error
	cmk := &cmk{}
	cmk.PublicKey, err = a.fetchPublicKey(ctx)
	if err != nil {
		return nil, err
	}
	cmk.KeyMetadata, err = a.fetchKeyMetadata(ctx)
	if err != nil {
		return nil, err
	}
	return cmk, nil
}

func (a *fakeAWSClient) getCMK(ctx context.Context) (*cmk, error) {
	var lerr error
	loader := ttlcache.LoaderFunc[string, cmk](
		func(c *ttlcache.Cache[string, cmk], key string) *ttlcache.Item[string, cmk] {
			var k *cmk
			k, lerr = a.fetchCMK(ctx)
			if lerr == nil {
				return c.Set(cacheKey, *k, time.Second*300)
			}
			return nil
		},
	)

	item := a.keyCache.Get(cacheKey, ttlcache.WithLoader[string, cmk](loader))
	if lerr == nil {
		cmk := item.Value()
		return &cmk, nil
	}
	return nil, lerr
}

// At the moment this function lies unused, but it is here for future if necessary

func (a *fakeAWSClient) verifyRemotely(ctx context.Context, sig, digest []byte) error {
	c, err := a.getCMK(ctx)
	if err != nil {
		return err
	}

	v, err := cryptoutil.NewVerifier(c.PublicKey, cryptoutil.VerifyWithHash(a.hash))
	if err != nil {
		return err
	}

	return v.Verify(bytes.NewReader(digest), sig)
}

func (a *fakeAWSClient) sign(ctx context.Context, digest []byte, _ crypto.Hash) ([]byte, error) {
	_, err := a.getCMK(ctx)
	if err != nil {
		return nil, err
	}

	signer, err := cryptoutil.NewSigner(a.privateKey, cryptoutil.SignWithHash(a.hash))
	if err != nil {
		return nil, err
	}

	s, err := signer.Sign(bytes.NewReader(digest))
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (a *fakeAWSClient) fetchPublicKey(ctx context.Context) (crypto.PublicKey, error) {
	k, p, err := createRsaKey()
	if err != nil {
		return nil, err
	}
	a.privateKey = k

	return p, nil
}

func (a *fakeAWSClient) fetchKeyMetadata(ctx context.Context) (*types.KeyMetadata, error) {
	km := &types.KeyMetadata{
		KeyId:        &a.keyID,
		AWSAccountId: &aid,
		Arn:          &arn,
	}

	return km, nil
}
