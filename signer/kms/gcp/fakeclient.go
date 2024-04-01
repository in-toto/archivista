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
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/kms/apiv1/kmspb"

	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/signer/kms"
	"github.com/jellydator/ttlcache/v3"
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

type fakeGCPClient struct {
	projectID  string
	locationID string
	keyRing    string
	keyName    string
	version    string
	kvCache    *ttlcache.Cache[string, cryptoKeyVersion]
	signer     cryptoutil.Signer
}

func newFakeGCPClient(ctx context.Context, ksp *kms.KMSSignerProvider) (*fakeGCPClient, error) {
	fmt.Println(ksp.Reference)
	if err := ValidReference(ksp.Reference); err != nil {
		return nil, err
	}

	g := &fakeGCPClient{
		kvCache: nil,
	}

	var err error
	g.projectID, g.locationID, g.keyRing, g.keyName, g.version, err = parseReference(ksp.Reference)
	if err != nil {
		return nil, err
	}

	g.kvCache = ttlcache.New[string, cryptoKeyVersion](
		ttlcache.WithDisableTouchOnHit[string, cryptoKeyVersion](),
	)

	// prime the cache
	g.kvCache.Get(cacheKey)
	return g, nil
}

func (g *fakeGCPClient) Verifier() (cryptoutil.Verifier, error) {
	crv, err := g.getCKV()
	if err != nil {
		return nil, fmt.Errorf("transient error while getting KMS verifier: %w", err)
	}

	return crv.Verifier, nil
}

// keyVersionName returns the first key version found for a key in KMS
func (g *fakeGCPClient) keyVersionName(ctx context.Context) (*cryptoKeyVersion, error) {
	pubKey, err := g.fetchPublicKey(ctx, "1")
	if err != nil {
		return nil, fmt.Errorf("unable to fetch public key while creating signer: %w", err)
	}

	// kv is keyVersion to use
	crv := cryptoKeyVersion{
		CryptoKeyVersion: &kmspb.CryptoKeyVersion{
			Name:      "1",
			Algorithm: kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_2048_SHA256,
		},
		PublicKey: pubKey,
	}

	g.version = "1"

	// crv.Verifier is set here to enable storing the public key & hash algorithm together,
	// as well as using the in memory Verifier to perform the verify operations.
	switch crv.CryptoKeyVersion.Algorithm {
	case kmspb.CryptoKeyVersion_EC_SIGN_P256_SHA256:
		pub, ok := pubKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("public key is not rsa")
		}
		crv.Verifier = cryptoutil.NewECDSAVerifier(pub, crypto.SHA256)
		crv.HashFunc = crypto.SHA256
	case kmspb.CryptoKeyVersion_EC_SIGN_P384_SHA384:
		pub, ok := pubKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("public key is not rsa")
		}
		crv.Verifier = cryptoutil.NewECDSAVerifier(pub, crypto.SHA384)
		crv.HashFunc = crypto.SHA384
	case kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_2048_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_3072_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_4096_SHA256:
		pub, ok := pubKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("public key is not rsa")
		}
		crv.Verifier = cryptoutil.NewRSAVerifier(pub, crypto.SHA256)
		crv.HashFunc = crypto.SHA256
	case kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_4096_SHA512:
		pub, ok := pubKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("public key is not rsa")
		}
		crv.Verifier = cryptoutil.NewRSAVerifier(pub, crypto.SHA384)
		crv.HashFunc = crypto.SHA384
	case kmspb.CryptoKeyVersion_RSA_SIGN_PSS_2048_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PSS_3072_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PSS_4096_SHA256:
		pub, ok := pubKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("public key is not rsa")
		}
		crv.Verifier = cryptoutil.NewRSAVerifier(pub, crypto.SHA256)
		crv.HashFunc = crypto.SHA256
	case kmspb.CryptoKeyVersion_RSA_SIGN_PSS_4096_SHA512:
		pub, ok := pubKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("public key is not rsa")
		}
		crv.Verifier = cryptoutil.NewRSAVerifier(pub, crypto.SHA512)
		crv.HashFunc = crypto.SHA512
	default:
		return nil, errors.New("unknown algorithm specified by KMS")
	}
	if err != nil {
		return nil, fmt.Errorf("initializing internal verifier: %w", err)
	}

	return &crv, nil
}

func (g *fakeGCPClient) fetchPublicKey(ctx context.Context, name string) (crypto.PublicKey, error) {
	priv, pub, err := createRsaKey()
	if err != nil {
		return nil, err
	}
	sign := cryptoutil.NewRSASigner(priv, crypto.SHA256)
	g.signer = sign

	return pub, nil
}

// getCKV gets the latest CryptoKeyVersion from the client's cache, which may trigger an actual
// call to GCP if the existing entry in the cache has expired.
func (g *fakeGCPClient) getCKV() (*cryptoKeyVersion, error) {
	var lerr error
	loader := ttlcache.LoaderFunc[string, cryptoKeyVersion](
		func(c *ttlcache.Cache[string, cryptoKeyVersion], key string) *ttlcache.Item[string, cryptoKeyVersion] {
			var ttl time.Duration
			var data *cryptoKeyVersion

			// if we're given an explicit version, cache this value forever
			if g.version != "" {
				ttl = time.Second * 0
			} else {
				ttl = time.Second * 300
			}
			data, lerr = g.keyVersionName(context.Background())
			if lerr == nil {
				return c.Set(key, *data, ttl)
			}
			return nil
		},
	)

	// we get once and use consistently to ensure the cache value doesn't change underneath us
	item := g.kvCache.Get(cacheKey, ttlcache.WithLoader[string, cryptoKeyVersion](loader))
	if item != nil {
		v := item.Value()
		return &v, nil
	}

	return nil, lerr
}

func (g *fakeGCPClient) sign(ctx context.Context, digest []byte, alg crypto.Hash, crc uint32) ([]byte, error) {
	_, err := g.getCKV()
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(digest)

	return g.signer.Sign(reader)
}

// Seems like GCP doesn't support any remote verification, so we'll just use the local verifier
func (g *fakeGCPClient) verify(message io.Reader, sig []byte) error {
	crv, err := g.getCKV()
	if err != nil {
		return fmt.Errorf("transient error getting info from KMS: %w", err)
	}

	if err := crv.Verifier.Verify(message, sig); err != nil {
		// key could have been rotated, clear cache and try again if we're not pinned to a version
		if g.version == "" {
			g.kvCache.Delete(cacheKey)
			crv, err = g.getCKV()
			if err != nil {
				return fmt.Errorf("transient error getting info from KMS: %w", err)
			}
			return crv.Verifier.Verify(message, sig)
		}
		return fmt.Errorf("failed to verify for fixed version: %w", err)
	}

	return nil
}
