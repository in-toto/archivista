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
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/in-toto/archivista/ent"
	"github.com/in-toto/archivista/ent/dsse"
	"github.com/in-toto/archivista/ent/sigstorebundle"
	witnessdsse "github.com/in-toto/go-witness/dsse"
)

// ReconstructBundleFromDSSE builds a minimal Sigstore bundle from stored DSSE data with provided payload
func ReconstructBundleFromDSSE(ctx context.Context, client *ent.Client, dsseID uuid.UUID, payload []byte) ([]byte, error) {
	// Fetch DSSE with all edges
	dsseEnt, err := client.Dsse.Query().
		Where(dsse.ID(dsseID)).
		WithSignatures(func(q *ent.SignatureQuery) {
			q.WithTimestamps()
		}).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	bundle := &Bundle{
		MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json",
	}

	// Reconstruct DSSE envelope with provided payload
	bundle.DsseEnvelope = &DsseEnvelope{
		Payload:     base64.StdEncoding.EncodeToString(payload),
		PayloadType: dsseEnt.PayloadType,
	}

	// Use first signature (Sigstore bundles prefer single sig)
	if len(dsseEnt.Edges.Signatures) > 0 {
		sig := dsseEnt.Edges.Signatures[0]

		sigBytes, _ := base64.StdEncoding.DecodeString(sig.Signature)
		bundle.DsseEnvelope.Signatures = []DsseSig{{
			KeyID: sig.KeyID,
			Sig:   base64.StdEncoding.EncodeToString(sigBytes),
		}}

		// Reconstruct VerificationMaterial from stored signature data
		bundle.VerificationMaterial = reconstructVerificationMaterial(sig)
	}

	return json.Marshal(bundle)
}

// reconstructVerificationMaterial rebuilds VerificationMaterial from stored signature/timestamp data
func reconstructVerificationMaterial(sig *ent.Signature) *VerificationMaterial {
	vm := &VerificationMaterial{}

	// Certificate chain
	if len(sig.Certificate) > 0 {
		if len(sig.Intermediates) > 0 {
			// Build chain: leaf + intermediates
			chain := &X509CertificateChain{
				Certificates: []Certificate{{
					RawBytes: base64.StdEncoding.EncodeToString(sig.Certificate),
				}},
			}

			for _, intermediate := range sig.Intermediates {
				chain.Certificates = append(chain.Certificates, Certificate{
					RawBytes: base64.StdEncoding.EncodeToString(intermediate),
				})
			}

			vm.X509CertificateChain = chain
		} else {
			// Standalone cert
			vm.Certificate = &Certificate{
				RawBytes: base64.StdEncoding.EncodeToString(sig.Certificate),
			}
		}
	}

	// RFC3161 timestamps
	if len(sig.Edges.Timestamps) > 0 {
		vm.TimestampVerificationData = &TimestampVerificationData{}

		for _, ts := range sig.Edges.Timestamps {
			if ts.Type == "tsp" && len(ts.Data) > 0 {
				vm.TimestampVerificationData.RFC3161Timestamps = append(
					vm.TimestampVerificationData.RFC3161Timestamps,
					RFC3161Timestamp{
						SignedTimestamp: base64.StdEncoding.EncodeToString(ts.Data),
					},
				)
			}
		}
	}

	return vm
}

// ExportGoWitnessDSSE exports a go-witness compatible DSSE envelope with all signatures and provided payload
func ExportGoWitnessDSSE(ctx context.Context, client *ent.Client, dsseID uuid.UUID, payload []byte) (*witnessdsse.Envelope, error) {
	// Fetch DSSE with all edges
	dsseEnt, err := client.Dsse.Query().
		Where(dsse.ID(dsseID)).
		WithSignatures(func(q *ent.SignatureQuery) {
			q.WithTimestamps()
		}).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	env := &witnessdsse.Envelope{
		Payload:     payload,
		PayloadType: dsseEnt.PayloadType,
	}

	// Include ALL signatures (go-witness supports multiple)
	for _, sig := range dsseEnt.Edges.Signatures {
		sigBytes, _ := base64.StdEncoding.DecodeString(sig.Signature)
		witnessSig := witnessdsse.Signature{
			KeyID:         sig.KeyID,
			Signature:     sigBytes,
			Certificate:   sig.Certificate,
			Intermediates: sig.Intermediates,
		}

		// Add timestamps
		for _, ts := range sig.Edges.Timestamps {
			if ts.Type == "tsp" {
				witnessSig.Timestamps = append(witnessSig.Timestamps,
					witnessdsse.SignatureTimestamp{
						Type: witnessdsse.TimestampRFC3161,
						Data: ts.Data,
					})
			}
		}

		env.Signatures = append(env.Signatures, witnessSig)
	}

	return env, nil
}

// ExportBundle retrieves a bundle for a given DSSE with provided payload
// If a bundle was stored with the DSSE, returns the gitoid and mediaType
// Otherwise reconstructs a minimal bundle from DSSE data
func ExportBundle(ctx context.Context, client *ent.Client, dsseID uuid.UUID, payload []byte) ([]byte, error) {
	// Check if a bundle is linked to this DSSE
	bundle, err := client.SigstoreBundle.Query().
		Where(sigstorebundle.HasDsseWith(dsse.ID(dsseID))).
		Only(ctx)

	if err == nil && bundle != nil {
		// Bundle exists - caller should fetch the raw bytes from blob storage
		return []byte(fmt.Sprintf(`{"gitoidSha256":"%s","mediaType":"%s"}`, bundle.GitoidSha256, bundle.MediaType)), nil
	}

	// No bundle stored, reconstruct from DSSE with provided payload
	return ReconstructBundleFromDSSE(ctx, client, dsseID, payload)
}
