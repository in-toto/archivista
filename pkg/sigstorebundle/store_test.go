// Copyright 2025 The Archivista Contributors
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

package sigstorebundle_test

import (
	"encoding/base64"
	"testing"

	"github.com/in-toto/archivista/pkg/sigstorebundle"
	"github.com/stretchr/testify/assert"
)

func TestParseBundleJSON(t *testing.T) {
	bundleJSON := `{
		"mediaType": "application/vnd.dev.sigstore.bundle.v0.3+json",
		"verificationMaterial": {
			"certificate": {
				"rawBytes": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJrVENDQVRlZ0F3SUJBZ0lSWWQycHJsRUFBQUJJQmlsYmlNeXZFZFRBS0JnZ3Foa2pPUFFRREJERElGQWpJCmFNUmt3RmdZRFZRUUxFdzkyYjI5elpVbHlkSFZsSUUxa2JHOGdSVEl3Q2dZRFZRUUdFd0pWVXpBZUZ3MHlNVEF4CkxqRTBNREkwT1RRNVZHRnBibmd4RmpBVUJnTlZCQUVUQTBKVlN6QmhCZ05WQkFNVGQzZDNjaTExZDIweVRXWjAKWVdkbEluSmxkSFZsTFVObGNuWmxjaTFSU1V3eEZEQVNCZ05WQkZSelpXNTBJRlZ1ZDJ4aGJtTmxNQjBHQTFVZApKUVFXTUJRR0NDc0dBUVFCZ2pWR01Bc0dBMVVkRHdRRUF3SUZvREFmQmdOVkhTTUVHREFXZ0JSdDAxRlZSekl1CjZTNzAxSjlLZXo0WWF6SlJJVEF3Q2dZSUtvWkl6ajBFQXdJRFJ3QXdSQUlnQW1SVE9QR0hxRkN4c2l1UEpQajQKQzlrRkd0RUtQMmJ2blV2TjA1STEyNFU9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
			}
		},
		"dsseEnvelope": {
			"payload": "eyJfdHlwZSI6Imh0dHBzOi8vaW4tdG90by5pby9zdGF0ZW1lbnQvdjEiLCJwcmVkaWNhdGVUeXBlIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9wcmVkaWNhdGUvdjEiLCJzdWJqZWN0cyI6W3sibmFtZSI6ImZpbGU6Ly9leGFtcGxlLnR4dCIsImRpZ2VzdHMiOnsiYWxnb3JpdGhtIjoic2hhMjU2IiwidmFsdWUiOiJhYmNkZWYifX1dfQ==",
			"payloadType": "application/vnd.in-toto+json",
			"signatures": [
				{
					"keyid": "test-key-id",
					"sig": "MEUCIQD8eXL5pmW3xY8L+pLKYr5YQp/8cYhcGb7XdxmPQYXQ5AIgHM1dXWHYsw9S3K5XqZ7X8xJ4mFVl2lzJ+HgC5vYLRkI="
				}
			]
		}
	}`

	bundle, err := sigstorebundle.ParseBundle([]byte(bundleJSON))
	assert.NoError(t, err)
	assert.NotNil(t, bundle)
	assert.Equal(t, "application/vnd.dev.sigstore.bundle.v0.3+json", bundle.MediaType)
	assert.NotNil(t, bundle.DsseEnvelope)
	assert.NotNil(t, bundle.VerificationMaterial)
}

func TestMapBundleToDSSE(t *testing.T) {
	bundleJSON := `{
		"mediaType": "application/vnd.dev.sigstore.bundle.v0.3+json",
		"verificationMaterial": {
			"certificate": {
				"rawBytes": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJrVENDQVRlZ0F3SUJBZ0lSWWQycHJsRUFBQUJJQmlsYmlNeXZFZFRBS0JnZ3Foa2pPUFFRREJERklGQWpJCmFNUmt3RmdZRFZRUUxFdzkyYjI5elpVbHlkSFZsSUUxa2JHOGdSVEl3Q2dZRFZRUUdFd0pWVXpBZUZ3MHlNVEF4CkxqRTBNREkwT1RRNVZHRnBibmd4RmpBVUJnTlZCQUVUQTBKVlN6QmhCZ05WQkFNVGQzZDNjaTExZDIweVRXWjAKWVdkbEluSmxkSFZsTFVObGNuWmxjaTFSU1V3eEZEQVNCZ05WQkZSelpXNTBJRlZ1ZDJ4aGJtTmxNQjBHQTFVZApKUVFXTUJRR0NDc0dBUVFCZ2pWR01Bc0dBMVVkRHdRRUF3SUZvREFmQmdOVkhTTUVHREFXZ0JSdDAxRlZSekl1CjZTNzAxSjlLZXo0WWF6SlJJVEF3Q2dZSUtvWkl6ajBFQXdJRFJ3QXdSQUlnQW1SVE9QR0hxRkN4c2l1UEpQajQKQzlrRkd0RUtQMmJ2blV2TjA1STEyNFU9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
			}
		},
		"dsseEnvelope": {
			"payload": "eyJfdHlwZSI6Imh0dHBzOi8vaW4tdG90by5pby9zdGF0ZW1lbnQvdjEiLCJwcmVkaWNhdGVUeXBlIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9wcmVkaWNhdGUvdjEiLCJzdWJqZWN0cyI6W3sibmFtZSI6ImZpbGU6Ly9leGFtcGxlLnR4dCIsImRpZ2VzdHMiOnsiYWxnb3JpdGhtIjoic2hhMjU2IiwidmFsdWUiOiJhYmNkZWYifX1dfQ==",
			"payloadType": "application/vnd.in-toto+json",
			"signatures": [
				{
					"keyid": "test-key-id",
					"sig": "MEUCIQD8eXL5pmW3xY8L+pLKYr5YQp/8cYhcGb7XdxmPQYXQ5AIgHM1dXWHYsw9S3K5XqZ7X8xJ4mFVl2lzJ+HgC5vYLRkI="
				}
			]
		}
	}`

	bundle, err := sigstorebundle.ParseBundle([]byte(bundleJSON))
	assert.NoError(t, err)

	env, err := sigstorebundle.MapBundleToDSSE(bundle)
	assert.NoError(t, err)
	assert.NotNil(t, env)
	assert.Equal(t, "application/vnd.in-toto+json", env.PayloadType)
	assert.Len(t, env.Signatures, 1)
	assert.Equal(t, "test-key-id", env.Signatures[0].KeyID)
	assert.NotEmpty(t, env.Signatures[0].Certificate)
}

func TestMapBundleToDSSEWithChain(t *testing.T) {
	// Create a minimal cert (just test that it gets parsed)
	certDER := []byte("CERT_DER_DATA")
	intDER := []byte("INT_DER_DATA")

	certB64 := base64.StdEncoding.EncodeToString(certDER)
	intB64 := base64.StdEncoding.EncodeToString(intDER)

	bundle := &sigstorebundle.Bundle{
		MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
		VerificationMaterial: &sigstorebundle.VerificationMaterial{
			X509CertificateChain: &sigstorebundle.X509CertificateChain{
				Certificates: []sigstorebundle.Certificate{
					{RawBytes: certB64},
					{RawBytes: intB64},
				},
			},
		},
		DsseEnvelope: &sigstorebundle.DsseEnvelope{
			Payload:     base64.StdEncoding.EncodeToString([]byte(`{"_type":"test"}`)),
			PayloadType: "application/vnd.in-toto+json",
			Signatures: []sigstorebundle.DsseSig{
				{
					KeyID: "key1",
					Sig:   base64.StdEncoding.EncodeToString([]byte("signature")),
				},
			},
		},
	}

	env, err := sigstorebundle.MapBundleToDSSE(bundle)
	assert.NoError(t, err)
	assert.Len(t, env.Signatures, 1)
	assert.Equal(t, certDER, env.Signatures[0].Certificate)
	assert.Len(t, env.Signatures[0].Intermediates, 1)
	assert.Equal(t, intDER, env.Signatures[0].Intermediates[0])
}

func TestMapBundleToDSSE_NilBundle(t *testing.T) {
	env, err := sigstorebundle.MapBundleToDSSE(nil)
	assert.Error(t, err)
	assert.Nil(t, env)
	assert.Contains(t, err.Error(), "bundle is nil")
}

func TestMapBundleToDSSE_MissingDsseEnvelope(t *testing.T) {
	bundle := &sigstorebundle.Bundle{
		MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
	}
	env, err := sigstorebundle.MapBundleToDSSE(bundle)
	assert.Error(t, err)
	assert.Nil(t, env)
	assert.Contains(t, err.Error(), "dsseEnvelope")
}

func TestMapBundleToDSSE_EmptyPayload(t *testing.T) {
	bundle := &sigstorebundle.Bundle{
		MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
		DsseEnvelope: &sigstorebundle.DsseEnvelope{
			Payload:     "",
			PayloadType: "application/vnd.in-toto+json",
		},
	}
	env, err := sigstorebundle.MapBundleToDSSE(bundle)
	assert.Error(t, err)
	assert.Nil(t, env)
	assert.Contains(t, err.Error(), "payload is empty")
}

func TestMapBundleToDSSE_InvalidBase64Payload(t *testing.T) {
	bundle := &sigstorebundle.Bundle{
		MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
		DsseEnvelope: &sigstorebundle.DsseEnvelope{
			Payload:     "!!!invalid-base64!!!",
			PayloadType: "application/vnd.in-toto+json",
		},
	}
	env, err := sigstorebundle.MapBundleToDSSE(bundle)
	assert.Error(t, err)
	assert.Nil(t, env)
	assert.Contains(t, err.Error(), "failed to decode dsseEnvelope.payload")
}

func TestMapBundleToDSSE_NoSignatures(t *testing.T) {
	bundle := &sigstorebundle.Bundle{
		MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
		DsseEnvelope: &sigstorebundle.DsseEnvelope{
			Payload:     base64.StdEncoding.EncodeToString([]byte(`{"_type":"test"}`)),
			PayloadType: "application/vnd.in-toto+json",
			Signatures:  []sigstorebundle.DsseSig{},
		},
	}
	env, err := sigstorebundle.MapBundleToDSSE(bundle)
	assert.Error(t, err)
	assert.Nil(t, env)
	assert.Contains(t, err.Error(), "no signatures")
}

func TestMapBundleToDSSE_MissingSignatureData(t *testing.T) {
	bundle := &sigstorebundle.Bundle{
		MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
		DsseEnvelope: &sigstorebundle.DsseEnvelope{
			Payload:     base64.StdEncoding.EncodeToString([]byte(`{"_type":"test"}`)),
			PayloadType: "application/vnd.in-toto+json",
			Signatures: []sigstorebundle.DsseSig{
				{
					KeyID: "key1",
					Sig:   "",
				},
			},
		},
	}
	env, err := sigstorebundle.MapBundleToDSSE(bundle)
	assert.Error(t, err)
	assert.Nil(t, env)
	assert.Contains(t, err.Error(), "signature 0 missing sig field")
}

func TestMapBundleToDSSE_InvalidBase64Signature(t *testing.T) {
	bundle := &sigstorebundle.Bundle{
		MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
		DsseEnvelope: &sigstorebundle.DsseEnvelope{
			Payload:     base64.StdEncoding.EncodeToString([]byte(`{"_type":"test"}`)),
			PayloadType: "application/vnd.in-toto+json",
			Signatures: []sigstorebundle.DsseSig{
				{
					KeyID: "key1",
					Sig:   "!!!invalid-base64!!!",
				},
			},
		},
	}
	env, err := sigstorebundle.MapBundleToDSSE(bundle)
	assert.Error(t, err)
	assert.Nil(t, env)
	assert.Contains(t, err.Error(), "failed to decode signature 0")
}
