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

// Bundle represents a Sigstore bundle v0.3+
type Bundle struct {
	MediaType            string                   `json:"mediaType"`
	VerificationMaterial *VerificationMaterial    `json:"verificationMaterial,omitempty"`
	MessageSignature     *MessageSignature        `json:"messageSignature,omitempty"`
	DsseEnvelope         *DsseEnvelope           `json:"dsseEnvelope,omitempty"`
}

// VerificationMaterial contains certificate chain and timestamps
type VerificationMaterial struct {
	X509CertificateChain      *X509CertificateChain      `json:"x509CertificateChain,omitempty"`
	Certificate               *Certificate                `json:"certificate,omitempty"`
	TimestampVerificationData *TimestampVerificationData `json:"timestampVerificationData,omitempty"`
}

// X509CertificateChain is a chain of X.509 certificates
type X509CertificateChain struct {
	Certificates []Certificate `json:"certificates"`
}

// Certificate is a base64-encoded DER certificate
type Certificate struct {
	RawBytes string `json:"rawBytes"` // base64 DER
}

// TimestampVerificationData contains RFC3161 timestamps
type TimestampVerificationData struct {
	RFC3161Timestamps []RFC3161Timestamp `json:"rfc3161Timestamps"`
}

// RFC3161Timestamp is a base64-encoded RFC3161 timestamp
type RFC3161Timestamp struct {
	SignedTimestamp string `json:"signedTimestamp"` // base64 DER
}

// DsseEnvelope from the bundle
type DsseEnvelope struct {
	Payload     string    `json:"payload"`     // base64
	PayloadType string    `json:"payloadType"`
	Signatures  []DsseSig `json:"signatures"`
}

// DsseSig is a single signature in a DSSE envelope
type DsseSig struct {
	KeyID string `json:"keyid"`
	Sig   string `json:"sig"` // base64
}

// MessageSignature from the bundle (not implemented in v0)
type MessageSignature struct {
	MessageDigest *MessageDigest `json:"messageDigest,omitempty"`
	Signature     string         `json:"signature"`
}

// MessageDigest for message signatures
type MessageDigest struct {
	Algorithm string `json:"algorithm"`
	Digest    string `json:"digest"`
}
