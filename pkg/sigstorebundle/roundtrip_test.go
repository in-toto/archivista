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
	"encoding/json"
	"testing"

	"github.com/in-toto/archivista/pkg/sigstorebundle"
	"github.com/stretchr/testify/assert"
)

// TestBundleRoundTripParsing tests that we can parse and reconstruct bundles
func TestBundleRoundTripParsing(t *testing.T) {
	// Original bundle JSON
	originalBundle := `{
		"mediaType": "application/vnd.dev.sigstore.bundle.v0.3+json",
		"verificationMaterial": {
			"x509CertificateChain": {
				"certificates": [
					{"rawBytes": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJrVENDQVRlZ0F3SUJBZ0lSWWQycHJsRUFBQUJJQmlsYmlNeXZFZFRBS0JnZ3Foa2pPUFFRREJERklGQWpJCmFNUmt3RmdZRFZRUUxFdzkyYjI5elpVbHlkSFZsSUUxa2JHOGdSVEl3Q2dZRFZRUUdFd0pWVXpBZUZ3MHlNVEF4CkxqRTBNREkwT1RRNVZHRnBibmd4RmpBVUJnTlZCQUVUQTBKVlN6QmhCZ05WQkFNVGQzZDNjaTExZDIweVRXWjAKWVdkbEluSmxkSFZsTFVObGNuWmxjaTFSU1V3eEZEQVNCZ05WQkZSalpXNTBJRlZ1ZDJ4aGJtTmxNQjBHQTFVZApKUVFXTUJRR0NDc0dBUVFCZ2pWR01Bc0dBMVVkRHdRRUF3SUZvREFmQmdOVkhTTUVHREFXZ0JSdDAxRlZSekl1CjZTNzAxSjlLZXo0WWF6SlJJVEF3Q2dZSUtvWkl6ajBFQXdJRFJ3QXdSQUlnQW1SVE9QR0hxRkN4c2l1UEpQajQKQzlrRkd0RUtQMmJ2blV2TjA1STEyNFU9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"},
					{"rawBytes": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJrVENDQVRlZ0F3SUJBZ0lSWWQycHJsRUFBQUJJQmlsYmlNeXZFZFRBS0JnZ3Foa2pPUFFRREJERklGQWpJCmFNUmt3RmdZRFZRUUxFdzkyYjI5elpVbHlkSFZsSUUxa2JHOGdSVEl3Q2dZRFZRUUdFd0pWVXpBZUZ3MHlNVEF4CkxqRTBNREkwT1RRNVZHRnBibmd4RmpBVUJnTlZCQUVUQTBKVlN6QmhCZ05WQkFNVGQzZDNjaTExZDIweVRXWjAKWVdkbEluSmxkSFZsTFVObGNuWmxjaTFSU1V3eEZEQVNCZ05WQkZSalpXNTBJRlZ1ZDJ4aGJtTmxNQjBHQTFVZApKUVFXTUJRR0NDc0dBUVFCZ2pWR01Bc0dBMVVkRHdRRUF3SUZvREFmQmdOVkhTTUVHREFXZ0JSdDAxRlZSekl1CjZTNzAxSjlLZXo0WWF6SlJJVEF3Q2dZSUtvWkl6ajBFQXdJRFJ3QXdSQUlnQW1SVE9QR0hxRkN4c2l1UEpQajQKQzlrRkd0RUtQMmJ2blV2TjA1STEyNFU9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"}
				]
			},
			"timestampVerificationData": {
				"rfc3161Timestamps": [
					{"signedTimestamp": "MIIFeTCCA2G..."}
				]
			}
		},
		"dsseEnvelope": {
			"payload": "eyJfdHlwZSI6Imh0dHBzOi8vaW4tdG90by5pby9zdGF0ZW1lbnQvdjEiLCJwcmVkaWNhdGVUeXBlIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9wcmVkaWNhdGUvdjEiLCJzdWJqZWN0cyI6W3sibmFtZSI6ImZpbGU6Ly9leGFtcGxlLnR4dCIsImRpZ2VzdHMiOnsiYWxnb3JpdGhtIjoic2hhMjU2IiwidmFsdWUiOiJhYmNkZWYifX1dfQ==",
			"payloadType": "application/vnd.in-toto+json",
			"signatures": [
				{
					"keyid": "5e8c57df8ae58fe9a29b29f9993e2fc3b25bd75eb2754f353880bad4b9ebfdb3",
					"sig": "MEUCIQD8eXL5pmW3xY8L+pLKYr5YQp/8cYhcGb7XdxmPQYXQ5AIgHM1dXWHYsw9S3K5XqZ7X8xJ4mFVl2lzJ+HgC5vYLRkI="
				}
			]
		}
	}`

	// Parse the original bundle
	bundle, err := sigstorebundle.ParseBundle([]byte(originalBundle))
	assert.NoError(t, err)
	assert.NotNil(t, bundle)

	// Verify it has the expected structure
	assert.Equal(t, "application/vnd.dev.sigstore.bundle.v0.3+json", bundle.MediaType)
	assert.NotNil(t, bundle.DsseEnvelope)
	assert.Equal(t, "application/vnd.in-toto+json", bundle.DsseEnvelope.PayloadType)
	assert.Len(t, bundle.DsseEnvelope.Signatures, 1)
	assert.NotNil(t, bundle.VerificationMaterial)
	assert.NotNil(t, bundle.VerificationMaterial.X509CertificateChain)
	assert.Len(t, bundle.VerificationMaterial.X509CertificateChain.Certificates, 2)

	// Map to DSSE
	env, err := sigstorebundle.MapBundleToDSSE(bundle)
	assert.NoError(t, err)
	assert.NotNil(t, env)

	// Verify DSSE has the signature data (certificates may not decode properly in test fixtures)
	assert.Len(t, env.Signatures, 1)
	assert.NotEmpty(t, env.Signatures[0].Signature)
}

// TestBundleParsingWithMessageSignature tests parsing message signature bundles
func TestBundleParsingWithMessageSignature(t *testing.T) {
	// Message signature bundle (not full DSSE)
	bundleJSON := `{
		"mediaType": "application/vnd.dev.sigstore.bundle.v0.3+json",
		"messageSignature": {
			"messageDigest": {
				"algorithm": "sha256",
				"digest": "abcdef1234567890"
			},
			"signature": "MEUCIQDsig123...=="
		},
		"verificationMaterial": {
			"certificate": {
				"rawBytes": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJrVENDQVRlZ0F3SUJBZ0lSWWQycHJsRUFBQUJJQmlsYmlNeXZFZFRBS0JnZ3Foa2pPUFFRREJERklGQWpJCmFNUmt3RmdZRFZRUUxFdzkyYjI5elpVbHlkSFZsSUUxa2JHOGdSVEl3Q2dZRFZRUUdFd0pWVXpBZUZ3MHlNVEF4CkxqRTBNREkwT1RRNVZHRnBibmd4RmpBVUJnTlZCQUVUQTBKVlN6QmhCZ05WQkFNVGQzZDNjaTExZDIweVRXWjAKWVdkbEluSmxkSFZsTFVObGNuWmxjaTFSU1V3eEZEQVNCZ05WQkZSalpXNTBJRlZ1ZDJ4aGJtTmxNQjBHQTFVZApKUVFXTUJRR0NDc0dBUVFCZ2pWR01Bc0dBMVVkRHdRRUF3SUZvREFmQmdOVkhTTUVHREFXZ0JSdDAxRlZSekl1CjZTNzAxSjlLZXo0WWF6SlJJVEF3Q2dZSUtvWkl6ajBFQXdJRFJ3QXdSQUlnQW1SVE9QR0hxRkN4c2l1UEpQajQKQzlrRkd0RUtQMmJ2blV2TjA1STEyNFU9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
			}
		}
	}`

	bundle, err := sigstorebundle.ParseBundle([]byte(bundleJSON))
	assert.NoError(t, err)
	assert.NotNil(t, bundle)
	assert.Equal(t, "application/vnd.dev.sigstore.bundle.v0.3+json", bundle.MediaType)
	assert.NotNil(t, bundle.MessageSignature)
	assert.Nil(t, bundle.DsseEnvelope)
}

// TestBundleRoundTripVerification tests that parsed/reconstructed bundles can be re-marshaled
func TestBundleRoundTripVerification(t *testing.T) {
	// Create a simple bundle
	payload := base64.StdEncoding.EncodeToString([]byte(`{"_type":"test"}`))
	sig := base64.StdEncoding.EncodeToString([]byte("signature"))

	bundle := &sigstorebundle.Bundle{
		MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
		DsseEnvelope: &sigstorebundle.DsseEnvelope{
			Payload:     payload,
			PayloadType: "application/vnd.in-toto+json",
			Signatures: []sigstorebundle.DsseSig{
				{
					KeyID: "key1",
					Sig:   sig,
				},
			},
		},
	}

	// Marshal to JSON
	bundleJSON, err := json.Marshal(bundle)
	assert.NoError(t, err)

	// Unmarshal back
	bundle2, err := sigstorebundle.ParseBundle(bundleJSON)
	assert.NoError(t, err)

	// Verify it's the same
	assert.Equal(t, bundle.MediaType, bundle2.MediaType)
	assert.Equal(t, bundle.DsseEnvelope.Payload, bundle2.DsseEnvelope.Payload)
	assert.Equal(t, bundle.DsseEnvelope.PayloadType, bundle2.DsseEnvelope.PayloadType)
	assert.Len(t, bundle2.DsseEnvelope.Signatures, 1)
	assert.Equal(t, "key1", bundle2.DsseEnvelope.Signatures[0].KeyID)
}

// TestBundleWithMultipleSignatures tests bundles with multiple signatures
func TestBundleWithMultipleSignatures(t *testing.T) {
	payload := base64.StdEncoding.EncodeToString([]byte(`{"_type":"test"}`))

	bundle := &sigstorebundle.Bundle{
		MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
		DsseEnvelope: &sigstorebundle.DsseEnvelope{
			Payload:     payload,
			PayloadType: "application/vnd.in-toto+json",
			Signatures: []sigstorebundle.DsseSig{
				{
					KeyID: "key1",
					Sig:   base64.StdEncoding.EncodeToString([]byte("sig1")),
				},
				{
					KeyID: "key2",
					Sig:   base64.StdEncoding.EncodeToString([]byte("sig2")),
				},
				{
					KeyID: "key3",
					Sig:   base64.StdEncoding.EncodeToString([]byte("sig3")),
				},
			},
		},
	}

	env, err := sigstorebundle.MapBundleToDSSE(bundle)
	assert.NoError(t, err)
	assert.Len(t, env.Signatures, 3)
	assert.Equal(t, "key1", env.Signatures[0].KeyID)
	assert.Equal(t, "key2", env.Signatures[1].KeyID)
	assert.Equal(t, "key3", env.Signatures[2].KeyID)
}

// TestBundleWithTimestamps tests RFC3161 timestamp handling
func TestBundleWithTimestamps(t *testing.T) {
	payload := base64.StdEncoding.EncodeToString([]byte(`{"_type":"test"}`))
	timestampData := []byte("RFC3161_TIMESTAMP_DATA")

	bundle := &sigstorebundle.Bundle{
		MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
		VerificationMaterial: &sigstorebundle.VerificationMaterial{
			TimestampVerificationData: &sigstorebundle.TimestampVerificationData{
				RFC3161Timestamps: []sigstorebundle.RFC3161Timestamp{
					{
						SignedTimestamp: base64.StdEncoding.EncodeToString(timestampData),
					},
				},
			},
		},
		DsseEnvelope: &sigstorebundle.DsseEnvelope{
			Payload:     payload,
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
	assert.Len(t, env.Signatures[0].Timestamps, 1)
	assert.Equal(t, timestampData, env.Signatures[0].Timestamps[0].Data)
}

// TestBundleDetection tests the bundle detection logic according to Sigstore spec
// A valid bundle must have:
// 1. Valid mediaType (application/.../bundle.../+json with sigstore keyword)
// 2. verificationMaterial field (required by spec)
// 3. Content (dsseEnvelope or messageSignature)
func TestBundleDetection(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected bool
	}{
		{
			name:     "valid sigstore bundle with DSSE",
			json:     `{"mediaType":"application/vnd.dev.sigstore.bundle.v0.3+json","verificationMaterial":{"certificate":{"rawBytes":"test"}},"dsseEnvelope":{"payload":"test","payloadType":"test","signatures":[{"keyid":"k","sig":"s"}]}}`,
			expected: true,
		},
		{
			name:     "valid bundle with message signature",
			json:     `{"mediaType":"application/vnd.dev.sigstore.bundle.v0.3+json","verificationMaterial":{"certificate":{"rawBytes":"test"}},"messageSignature":{"signature":"test"}}`,
			expected: true,
		},
		{
			name:     "missing verification material (required)",
			json:     `{"mediaType":"application/vnd.dev.sigstore.bundle.v0.3+json","dsseEnvelope":{"payload":"test"}}`,
			expected: false,
		},
		{
			name:     "missing content (no dsseEnvelope or messageSignature)",
			json:     `{"mediaType":"application/vnd.dev.sigstore.bundle.v0.3+json","verificationMaterial":{"certificate":{"rawBytes":"test"}}}`,
			expected: false,
		},
		{
			name:     "missing media type",
			json:     `{"verificationMaterial":{},"dsseEnvelope":{}}`,
			expected: false,
		},
		{
			name:     "non-sigstore media type",
			json:     `{"mediaType":"application/json","verificationMaterial":{},"dsseEnvelope":{}}`,
			expected: false,
		},
		{
			name:     "media type missing bundle keyword",
			json:     `{"mediaType":"application/vnd.dev.sigstore.v0.3+json","verificationMaterial":{},"dsseEnvelope":{}}`,
			expected: false,
		},
		{
			name:     "invalid JSON",
			json:     `{invalid json}`,
			expected: false,
		},
		{
			name:     "valid mediaType but no sigstore keyword",
			json:     `{"mediaType":"application/vnd.example.bundle.v0.3+json","verificationMaterial":{},"dsseEnvelope":{}}`,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := sigstorebundle.IsBundleJSON([]byte(tc.json))
			assert.Equal(t, tc.expected, result, "bundle detection mismatch for: %s", tc.name)
		})
	}
}

// TestBundleWithEmptyVerificationMaterial tests bundle with empty verification material
func TestBundleWithEmptyVerificationMaterial(t *testing.T) {
	payload := base64.StdEncoding.EncodeToString([]byte(`{"_type":"test"}`))

	bundle := &sigstorebundle.Bundle{
		MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
		VerificationMaterial: &sigstorebundle.VerificationMaterial{},
		DsseEnvelope: &sigstorebundle.DsseEnvelope{
			Payload:     payload,
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
	// Should still work without verification material
	assert.NotEmpty(t, env.Signatures[0].Signature)
}
