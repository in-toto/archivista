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
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/in-toto/archivista/ent"
	"github.com/in-toto/go-witness/dsse"
)

const (
	// maxPayloadSize is the maximum size of a decoded payload in bytes (100MB)
	maxPayloadSize = 100 * 1024 * 1024
	// maxSignaturesPerBundle is the maximum number of signatures allowed per bundle
	maxSignaturesPerBundle = 100
)

// gitoidSHA256 computes gitoid v1 SHA256 hash (blob header + content)
// Note: Currently unused but kept for future gitoid calculation needs
// nolint:unused
func gitoidSHA256(content []byte) string {
	h := sha256.New()
	fmt.Fprintf(h, "blob %d\x00", len(content))
	h.Write(content)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ParseBundle unmarshals and validates a Sigstore bundle
func ParseBundle(raw []byte) (*Bundle, error) {
	var b Bundle
	if err := json.Unmarshal(raw, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// MapBundleToDSSE converts a Sigstore bundle to a go-witness DSSE envelope
func MapBundleToDSSE(bundle *Bundle) (*dsse.Envelope, error) {
	if bundle == nil {
		return nil, fmt.Errorf("bundle is nil")
	}

	if bundle.DsseEnvelope == nil {
		return nil, fmt.Errorf("bundle missing required field: dsseEnvelope")
	}

	if bundle.DsseEnvelope.Payload == "" {
		return nil, fmt.Errorf("dsseEnvelope.payload is empty")
	}

	// Check payload size before decoding (base64 encoded size * 3/4 ≈ decoded size)
	estimatedSize := len(bundle.DsseEnvelope.Payload) * 3 / 4
	if estimatedSize > maxPayloadSize {
		return nil, fmt.Errorf("payload size (%d bytes) exceeds maximum allowed size (%d bytes)", estimatedSize, maxPayloadSize)
	}

	// Decode payload
	payload, err := base64.StdEncoding.DecodeString(bundle.DsseEnvelope.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode dsseEnvelope.payload: %w", err)
	}

	if len(payload) == 0 {
		return nil, fmt.Errorf("decoded payload is empty")
	}

	env := &dsse.Envelope{
		Payload:     payload,
		PayloadType: bundle.DsseEnvelope.PayloadType,
	}

	if len(bundle.DsseEnvelope.Signatures) == 0 {
		return nil, fmt.Errorf("bundle has no signatures")
	}

	// Check signature count limit to prevent resource exhaustion
	if len(bundle.DsseEnvelope.Signatures) > maxSignaturesPerBundle {
		return nil, fmt.Errorf("bundle has %d signatures, exceeds maximum allowed (%d)", len(bundle.DsseEnvelope.Signatures), maxSignaturesPerBundle)
	}

	// Map signatures with VerificationMaterial
	for idx, bundleSig := range bundle.DsseEnvelope.Signatures {
		if bundleSig.Sig == "" {
			return nil, fmt.Errorf("signature %d missing sig field", idx)
		}

		sig, err := base64.StdEncoding.DecodeString(bundleSig.Sig)
		if err != nil {
			return nil, fmt.Errorf("failed to decode signature %d: %w", idx, err)
		}

		if len(sig) == 0 {
			return nil, fmt.Errorf("decoded signature %d is empty", idx)
		}

		witnessSig := dsse.Signature{
			KeyID:     bundleSig.KeyID,
			Signature: sig,
		}

		// Map VerificationMaterial to signature fields
		if bundle.VerificationMaterial != nil {
			mapVerificationMaterial(&witnessSig, bundle.VerificationMaterial)
		}

		env.Signatures = append(env.Signatures, witnessSig)
	}

	return env, nil
}

// mapVerificationMaterial maps Sigstore VerificationMaterial to go-witness signature fields
func mapVerificationMaterial(sig *dsse.Signature, vm *VerificationMaterial) {
	if sig == nil || vm == nil {
		return
	}

	// Certificate chain
	if vm.X509CertificateChain != nil && len(vm.X509CertificateChain.Certificates) > 0 {
		// First cert is leaf
		if vm.X509CertificateChain.Certificates[0].RawBytes != "" {
			if cert, err := base64.StdEncoding.DecodeString(vm.X509CertificateChain.Certificates[0].RawBytes); err == nil {
				sig.Certificate = cert
			}
		}

		// Rest are intermediates
		for i := 1; i < len(vm.X509CertificateChain.Certificates); i++ {
			if vm.X509CertificateChain.Certificates[i].RawBytes != "" {
				if cert, err := base64.StdEncoding.DecodeString(vm.X509CertificateChain.Certificates[i].RawBytes); err == nil {
					sig.Intermediates = append(sig.Intermediates, cert)
				}
			}
		}
	} else if vm.Certificate != nil && vm.Certificate.RawBytes != "" {
		// Standalone certificate
		if cert, err := base64.StdEncoding.DecodeString(vm.Certificate.RawBytes); err == nil {
			sig.Certificate = cert
		}
	}

	// RFC3161 timestamps
	if vm.TimestampVerificationData != nil {
		for _, ts := range vm.TimestampVerificationData.RFC3161Timestamps {
			if ts.SignedTimestamp != "" {
				if data, err := base64.StdEncoding.DecodeString(ts.SignedTimestamp); err == nil {
					sig.Timestamps = append(sig.Timestamps, dsse.SignatureTimestamp{
						Type: dsse.TimestampRFC3161,
						Data: data,
					})
				}
			}
		}
	}
}

// IsBundleJSON checks if the raw JSON represents a Sigstore bundle according to the
// official Sigstore bundle specification (https://github.com/sigstore/protobuf-specs)
//
// A valid bundle must have:
// 1. A mediaType matching the pattern: application/vnd.dev.sigstore.bundle.v<X>+json
// 2. A verificationMaterial field with key material
// 3. Content (either dsseEnvelope or messageSignature)
func IsBundleJSON(obj []byte) bool {
	// Try to unmarshal as a Bundle to validate structure
	bundle := &Bundle{}
	if err := json.Unmarshal(obj, bundle); err != nil {
		return false
	}

	// Validate mediaType according to spec: application/vnd.dev.sigstore.bundle.v<X>+json
	if bundle.MediaType == "" {
		return false
	}

	// Must match the official bundle media type pattern
	if !isValidBundleMediaType(bundle.MediaType) {
		return false
	}

	// According to protobuf spec, Bundle must have verificationMaterial
	if bundle.VerificationMaterial == nil {
		return false
	}

	// Must have at least one content type (DSSE or message signature)
	if bundle.DsseEnvelope == nil && bundle.MessageSignature == nil {
		return false
	}

	return true
}

// isValidBundleMediaType validates the mediaType according to Sigstore spec
// Valid formats: application/vnd.dev.sigstore.bundle.v<VERSION>+json
func isValidBundleMediaType(mediaType string) bool {
	// Check prefix and suffix
	if !strings.HasPrefix(mediaType, "application/") {
		return false
	}
	if !strings.HasSuffix(mediaType, "+json") {
		return false
	}

	// Must contain bundle and sigstore keywords
	if !strings.Contains(mediaType, "bundle") || !strings.Contains(mediaType, "sigstore") {
		return false
	}

	// Should match pattern: application/vnd.dev.sigstore.bundle.v<VERSION>+json
	// or: application/dev.sigstore.bundle+json (older format)
	// Accept flexible variations for compatibility
	return strings.Contains(mediaType, "bundle") &&
		(strings.Contains(mediaType, "sigstore") || strings.Contains(mediaType, "dev.sigstore"))
}

// StoreBundleMetadata stores bundle metadata in the database
// The raw bundle is stored in blob storage by the caller
func StoreBundleMetadata(ctx context.Context, tx *ent.Tx, gitoid string, mediaType string, dsseID uuid.UUID) (*ent.SigstoreBundle, error) {
	version := "0.3" // Parse from mediaType in future versions

	builder := tx.SigstoreBundle.Create().
		SetGitoidSha256(gitoid).
		SetMediaType(mediaType).
		SetVersion(version)

	// Link to DSSE if provided
	if dsseID != uuid.Nil {
		builder = builder.SetDsseID(dsseID)
	}

	bundle, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	return bundle, nil
}
