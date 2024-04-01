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
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"regexp"
	"strings"
	"time"

	gcpkms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/registry"
	"github.com/in-toto/go-witness/signer"
	"github.com/in-toto/go-witness/signer/kms"
	"github.com/jellydator/ttlcache/v3"
)

func init() {
	kms.AddProvider(ReferenceScheme, &gcpClientOptions{}, func(ctx context.Context, ksp *kms.KMSSignerProvider) (cryptoutil.Signer, error) {
		return LoadSignerVerifier(ctx, ksp)
	})
}

//nolint:revive
const (
	AlgorithmECDSAP256SHA256       = "ecdsa-p256-sha256"
	AlgorithmECDSAP384SHA384       = "ecdsa-p384-sha384"
	AlgorithmRSAPKCS1v152048SHA256 = "rsa-pkcs1v15-2048-sha256"
	AlgorithmRSAPKCS1v153072SHA256 = "rsa-pkcs1v15-3072-sha256"
	AlgorithmRSAPKCS1v154096SHA256 = "rsa-pkcs1v15-4096-sha256"
	AlgorithmRSAPKCS1v154096SHA512 = "rsa-pkcs1v15-4096-sha512"
	AlgorithmRSAPSS2048SHA256      = "rsa-pss-2048-sha256"
	AlgorithmRSAPSS3072SHA256      = "rsa-pss-3072-sha256"
	AlgorithmRSAPSS4096SHA256      = "rsa-pss-4096-sha256"
	AlgorithmRSAPSS4096SHA512      = "rsa-pss-4096-sha512"
)

var algorithmMap = map[string]kmspb.CryptoKeyVersion_CryptoKeyVersionAlgorithm{
	AlgorithmECDSAP256SHA256:       kmspb.CryptoKeyVersion_EC_SIGN_P256_SHA256,
	AlgorithmECDSAP384SHA384:       kmspb.CryptoKeyVersion_EC_SIGN_P384_SHA384,
	AlgorithmRSAPKCS1v152048SHA256: kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_2048_SHA256,
	AlgorithmRSAPKCS1v153072SHA256: kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_3072_SHA256,
	AlgorithmRSAPKCS1v154096SHA256: kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_4096_SHA256,
	AlgorithmRSAPKCS1v154096SHA512: kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_4096_SHA512,
	AlgorithmRSAPSS2048SHA256:      kmspb.CryptoKeyVersion_RSA_SIGN_PSS_2048_SHA256,
	AlgorithmRSAPSS3072SHA256:      kmspb.CryptoKeyVersion_RSA_SIGN_PSS_3072_SHA256,
	AlgorithmRSAPSS4096SHA256:      kmspb.CryptoKeyVersion_RSA_SIGN_PSS_4096_SHA256,
	AlgorithmRSAPSS4096SHA512:      kmspb.CryptoKeyVersion_RSA_SIGN_PSS_4096_SHA512,
}

var (
	errKMSReference = errors.New("kms specification should be in the format gcpkms://projects/[PROJECT_ID]/locations/[LOCATION]/keyRings/[KEY_RING]/cryptoKeys/[KEY]/cryptoKeyVersions/[VERSION]")

	re = regexp.MustCompile(`^gcpkms://projects/([^/]+)/locations/([^/]+)/keyRings/([^/]+)/cryptoKeys/([^/]+)(?:/(?:cryptoKeyVersions|versions)/([^/]+))?$`)
)

// ReferenceScheme schemes for various KMS services are copied from https://github.com/google/go-cloud/tree/master/secrets
const ReferenceScheme = "gcpkms://"

// ValidReference returns a non-nil error if the reference string is invalid
func ValidReference(ref string) error {
	if !re.MatchString(ref) {
		return errKMSReference
	}
	return nil
}

func parseReference(resourceID string) (projectID, locationID, keyRing, keyName, version string, err error) {
	v := re.FindStringSubmatch(resourceID)
	if len(v) != 6 {
		err = fmt.Errorf("invalid gcpkms format %q", resourceID)
		return
	}
	projectID, locationID, keyRing, keyName, version = v[1], v[2], v[3], v[4], v[5]
	return
}

type gcpClient struct {
	projectID  string
	locationID string
	keyRing    string
	keyName    string
	version    string
	kvCache    *ttlcache.Cache[string, cryptoKeyVersion]
	client     *gcpkms.KeyManagementClient
	options    *gcpClientOptions
}

type gcpClientOptions struct {
	credentialsFile string
}

type Option func(*gcpClientOptions)

func (a *gcpClientOptions) Init() []registry.Configurer {
	return []registry.Configurer{
		registry.StringConfigOption(
			"credentials-file",
			"The credentials file to use with the GCP KMS signer provider",
			"",
			func(sp signer.SignerProvider, cred string) (signer.SignerProvider, error) {
				ksp, ok := sp.(*kms.KMSSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a kms signer provider")
				}

				var clientOpts *gcpClientOptions
				for _, opt := range ksp.Options {
					co, optsOk := opt.(*gcpClientOptions)
					if !optsOk {
						continue
					}
					clientOpts = co
				}

				if clientOpts == nil {
					return nil, fmt.Errorf("unable to find aws client options in aws kms signer provider")
				}

				WithCredentialsFile(cred)(clientOpts)
				return ksp, nil
			},
		),
	}
}

func (*gcpClientOptions) ProviderName() string {
	name := fmt.Sprintf("kms-%s", strings.TrimSuffix(ReferenceScheme, "kms://"))
	return name
}

func WithCredentialsFile(cred string) Option {
	return func(opts *gcpClientOptions) {
		opts.credentialsFile = cred
	}
}

func newGCPClient(ctx context.Context, ksp *kms.KMSSignerProvider) (*gcpClient, error) {
	if err := ValidReference(ksp.Reference); err != nil {
		return nil, err
	}

	if ctx == nil {
		ctx = context.TODO()
	}

	g := &gcpClient{
		kvCache: nil,
	}

	var err error
	g.projectID, g.locationID, g.keyRing, g.keyName, g.version, err = parseReference(ksp.Reference)
	if err != nil {
		return nil, err
	}

	var ok bool
	for _, opt := range ksp.Options {
		g.options, ok = opt.(*gcpClientOptions)
		if ok {
			break
		}
	}

	if g.options == nil {
		return nil, fmt.Errorf("unable to find gcp client options in gcp kms signer provider")
	}

	var opts []option.ClientOption
	if g.options.credentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(g.options.credentialsFile))
	}

	g.client, err = gcpkms.NewKeyManagementClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("new gcp kms client: %w", err)
	}

	g.kvCache = ttlcache.New[string, cryptoKeyVersion](
		ttlcache.WithDisableTouchOnHit[string, cryptoKeyVersion](),
	)

	// prime the cache
	g.kvCache.Get(cacheKey)
	return g, nil
}

type cryptoKeyVersion struct {
	CryptoKeyVersion *kmspb.CryptoKeyVersion
	Verifier         cryptoutil.Verifier
	PublicKey        crypto.PublicKey
	HashFunc         crypto.Hash
}

// use a consistent key for cache lookups
const cacheKey = "crypto_key_version"

func (g *gcpClient) Verifier() (cryptoutil.Verifier, error) {
	crv, err := g.getCKV()
	if err != nil {
		return nil, fmt.Errorf("transient error while getting KMS verifier: %w", err)
	}

	return crv.Verifier, nil
}

// keyVersionName returns the first key version found for a key in KMS
func (g *gcpClient) keyVersionName(ctx context.Context) (*cryptoKeyVersion, error) {
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", g.projectID, g.locationID, g.keyRing, g.keyName)

	parentReq := &kmspb.GetCryptoKeyRequest{
		Name: parent,
	}
	key, err := g.client.GetCryptoKey(ctx, parentReq)
	if err != nil {
		return nil, err
	}
	if key.Purpose != kmspb.CryptoKey_ASYMMETRIC_SIGN {
		return nil, errors.New("specified key cannot be used to sign")
	}

	// if g.version was specified, use it explicitly
	var kv *kmspb.CryptoKeyVersion
	if g.version != "" {
		req := &kmspb.GetCryptoKeyVersionRequest{
			Name: parent + fmt.Sprintf("/cryptoKeyVersions/%s", g.version),
		}
		kv, err = g.client.GetCryptoKeyVersion(ctx, req)
		if err != nil {
			return nil, err
		}
	} else {
		req := &kmspb.ListCryptoKeyVersionsRequest{
			Parent:  parent,
			Filter:  "state=ENABLED",
			OrderBy: "name desc",
		}
		iterator := g.client.ListCryptoKeyVersions(ctx, req)

		// pick the key version that is enabled with the greatest version value
		kv, err = iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("unable to find an enabled key version in GCP KMS: %w", err)
		}
	}

	pubKey, err := g.fetchPublicKey(ctx, kv.Name)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch public key while creating signer: %w", err)
	}

	// kv is keyVersion to use
	crv := cryptoKeyVersion{
		CryptoKeyVersion: kv,
		PublicKey:        pubKey,
	}

	switch kv.Algorithm {
	case kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_2048_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_3072_SHA256,
		kmspb.CryptoKeyVersion_EC_SIGN_P256_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PSS_2048_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PSS_3072_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PSS_4096_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_4096_SHA256:
		crv.HashFunc = crypto.SHA256
	case kmspb.CryptoKeyVersion_EC_SIGN_P384_SHA384:
		crv.HashFunc = crypto.SHA384
	case kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_4096_SHA512,
		kmspb.CryptoKeyVersion_RSA_SIGN_PSS_4096_SHA512:
		crv.HashFunc = crypto.SHA512
	default:
		return nil, errors.New("unknown algorithm specified by KMS")
	}

	crv.Verifier, err = cryptoutil.NewVerifier(pubKey, cryptoutil.VerifyWithHash(crv.HashFunc))
	if err != nil {
		return nil, fmt.Errorf("initializing internal verifier: %w", err)
	}

	return &crv, nil
}

func (g *gcpClient) fetchPublicKey(ctx context.Context, name string) (crypto.PublicKey, error) {
	// Build the request.
	pkreq := &kmspb.GetPublicKeyRequest{Name: name}
	// Call the API.
	pk, err := g.client.GetPublicKey(ctx, pkreq)
	if err != nil {
		return nil, fmt.Errorf("public key: %w", err)
	}
	return cryptoutil.UnmarshalPEMToPublicKey([]byte(pk.GetPem()))
}

func (g *gcpClient) getHashFunc() (crypto.Hash, error) {
	ckv, err := g.getCKV()
	if err != nil {
		return 0, err
	}
	return ckv.HashFunc, nil
}

// getCKV gets the latest CryptoKeyVersion from the client's cache, which may trigger an actual
// call to GCP if the existing entry in the cache has expired.
func (g *gcpClient) getCKV() (*cryptoKeyVersion, error) {
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
			data, lerr = g.keyVersionName(context.TODO())
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

func (g *gcpClient) sign(ctx context.Context, digest []byte, alg crypto.Hash, crc uint32) ([]byte, error) {
	ckv, err := g.getCKV()
	if err != nil {
		return nil, err
	}

	gcpSignReq := kmspb.AsymmetricSignRequest{
		Name:   ckv.CryptoKeyVersion.Name,
		Digest: &kmspb.Digest{},
	}

	if crc != 0 {
		gcpSignReq.DigestCrc32C = wrapperspb.Int64(int64(crc))
	}

	switch alg {
	case crypto.SHA256:
		gcpSignReq.Digest.Digest = &kmspb.Digest_Sha256{
			Sha256: digest,
		}
	case crypto.SHA384:
		gcpSignReq.Digest.Digest = &kmspb.Digest_Sha384{
			Sha384: digest,
		}
	case crypto.SHA512:
		gcpSignReq.Digest.Digest = &kmspb.Digest_Sha512{
			Sha512: digest,
		}
	default:
		return nil, errors.New("unsupported hash function")
	}

	resp, err := g.client.AsymmetricSign(ctx, &gcpSignReq)
	if err != nil {
		return nil, fmt.Errorf("calling GCP AsymmetricSign: %w", err)
	}

	// Optional, but recommended: perform integrity verification on result.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	if crc != 0 && !resp.VerifiedDigestCrc32C {
		return nil, fmt.Errorf("AsymmetricSign: request corrupted in-transit")
	}
	if int64(crc32.Checksum(resp.Signature, crc32.MakeTable(crc32.Castagnoli))) != resp.SignatureCrc32C.Value {
		return nil, fmt.Errorf("AsymmetricSign: response corrupted in-transit")
	}

	return resp.Signature, nil
}

func (g *gcpClient) public(ctx context.Context) (crypto.PublicKey, error) {
	crv, err := g.getCKV()
	if err != nil {
		return nil, fmt.Errorf("transient error getting info from KMS: %w", err)
	}
	return crv.PublicKey, nil
}

// Seems like GCP doesn't support any remote verification, so we'll just use the local verifier
func (g *gcpClient) verify(message io.Reader, sig []byte) error {
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
