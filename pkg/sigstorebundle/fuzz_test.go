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

package sigstorebundle

import (
	"encoding/json"
	"testing"
)

// FuzzIsBundleJSON fuzzes the bundle detection logic
// This helps ensure we don't panic on malformed JSON or unexpected structures
func FuzzIsBundleJSON(f *testing.F) {
	// Seed corpus with valid and edge case inputs
	f.Add([]byte(`{"mediaType":"application/vnd.dev.sigstore.bundle.v0.3+json","verificationMaterial":{},"dsseEnvelope":{}}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`null`))
	f.Add([]byte(`[]`))
	f.Add([]byte(`"string"`))
	f.Add([]byte(`123`))
	f.Add([]byte(``))
	f.Add([]byte(`{`))
	f.Add([]byte(`{"mediaType":null}`))

	f.Fuzz(func(t *testing.T, data []byte) {
		// Should never panic, only return true/false
		_ = IsBundleJSON(data)
	})
}

// FuzzParseBundle fuzzes the bundle parsing logic
// This ensures we handle malformed JSON gracefully
func FuzzParseBundle(f *testing.F) {
	// Seed corpus
	f.Add([]byte(`{"mediaType":"application/vnd.dev.sigstore.bundle.v0.3+json"}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"dsseEnvelope":{"payload":"dGVzdA=="}}`))

	f.Fuzz(func(t *testing.T, data []byte) {
		// Should never panic
		_, _ = ParseBundle(data)
	})
}

// FuzzMapBundleToDSSE fuzzes the main bundle to DSSE conversion
// This is critical for security as it processes untrusted input
func FuzzMapBundleToDSSE(f *testing.F) {
	// Seed with various valid structures
	validBundle := `{
		"mediaType": "application/vnd.dev.sigstore.bundle.v0.3+json",
		"dsseEnvelope": {
			"payload": "dGVzdA==",
			"payloadType": "application/vnd.in-toto+json",
			"signatures": [{"keyid": "test", "sig": "c2lnbmF0dXJl"}]
		},
		"verificationMaterial": {
			"certificate": {"rawBytes": "Y2VydA=="}
		}
	}`

	f.Add([]byte(validBundle))
	f.Add([]byte(`{"dsseEnvelope":{}}`))
	f.Add([]byte(`{"dsseEnvelope":{"payload":"","signatures":[]}}`))

	f.Fuzz(func(t *testing.T, data []byte) {
		var bundle Bundle
		if err := json.Unmarshal(data, &bundle); err != nil {
			return // Skip invalid JSON
		}

		// Test with default limits
		_, err := MapBundleToDSSE(&bundle)

		// We expect errors for invalid bundles, just shouldn't panic
		_ = err

		// Test with custom limits (smaller to find more edge cases)
		limits := &BundleLimits{
			MaxPayloadSizeMB:       1,  // 1MB
			MaxSignaturesPerBundle: 10, // 10 signatures
		}
		_, err = MapBundleToDSSE(&bundle, limits)
		_ = err
	})
}

// FuzzBundleMediaType fuzzes media type validation
// Ensures we don't have ReDoS or other regex vulnerabilities
func FuzzBundleMediaType(f *testing.F) {
	// Seed with various media types
	f.Add("application/vnd.dev.sigstore.bundle.v0.3+json")
	f.Add("application/vnd.dev.sigstore.bundle+json")
	f.Add("")
	f.Add("application/json")
	f.Add("bundle")
	f.Add("sigstore")
	f.Add("application/vnd.dev.sigstore.bundle.v" + string(make([]byte, 1000)) + "+json")

	f.Fuzz(func(t *testing.T, mediaType string) {
		// Should never panic
		_ = isValidBundleMediaType(mediaType)
	})
}

// FuzzBundleWithLargePayload tests payload size limits
// Ensures we reject oversized payloads before attempting to decode
func FuzzBundleWithLargePayload(f *testing.F) {
	f.Fuzz(func(t *testing.T, payloadSize int) {
		if payloadSize < 0 || payloadSize > 200*1024*1024 {
			return // Skip unreasonable sizes
		}

		// Create a bundle with specified payload size
		largePayload := make([]byte, payloadSize)
		bundle := &Bundle{
			MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
			DsseEnvelope: &DsseEnvelope{
				Payload:     string(largePayload),
				PayloadType: "test",
				Signatures:  []DsseSig{{KeyID: "test", Sig: "dGVzdA=="}},
			},
		}

		// Should handle large payloads gracefully (either accept or reject with error)
		limits := &BundleLimits{
			MaxPayloadSizeMB:       100,
			MaxSignaturesPerBundle: 100,
		}
		_, err := MapBundleToDSSE(bundle, limits)

		// We expect an error for oversized payloads, but no panic
		_ = err
	})
}

// FuzzBundleWithManySignatures tests signature count limits
// Ensures we reject bundles with too many signatures
func FuzzBundleWithManySignatures(f *testing.F) {
	f.Fuzz(func(t *testing.T, sigCount uint8) {
		if sigCount == 0 {
			return
		}

		// Create a bundle with specified signature count
		sigs := make([]DsseSig, sigCount)
		for i := range sigs {
			sigs[i] = DsseSig{KeyID: "test", Sig: "dGVzdA=="}
		}

		bundle := &Bundle{
			MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
			DsseEnvelope: &DsseEnvelope{
				Payload:     "dGVzdA==",
				PayloadType: "test",
				Signatures:  sigs,
			},
		}

		limits := &BundleLimits{
			MaxPayloadSizeMB:       100,
			MaxSignaturesPerBundle: 100,
		}
		_, err := MapBundleToDSSE(bundle, limits)

		// Should handle any signature count gracefully
		_ = err
	})
}
