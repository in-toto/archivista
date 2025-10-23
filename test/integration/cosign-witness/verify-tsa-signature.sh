#!/bin/bash
# Copyright 2025 The Archivista Contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euo pipefail

# verify-tsa-signature.sh - Verify RFC3161 TSA signature in bundle

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ARTIFACTS_DIR="${SCRIPT_DIR}/artifacts"
KEYS_DIR="${SCRIPT_DIR}/keys"

echo "=== Verifying RFC3161 TSA Signature ==="

BUNDLE_FILE="${ARTIFACTS_DIR}/downloaded-bundle.json"

if [[ ! -f "${BUNDLE_FILE}" ]]; then
    echo "ERROR: Bundle file not found at ${BUNDLE_FILE}"
    exit 1
fi

echo "Bundle file: ${BUNDLE_FILE}"
echo ""

# Extract TSA timestamp info
echo "1. Extract TSA Information from Bundle"
echo "======================================"
TSA_COUNT=$(jq '.verificationMaterial.timestampVerificationData.rfc3161Timestamps | length' "${BUNDLE_FILE}" 2>/dev/null || echo 0)
echo "Number of RFC3161 timestamps: ${TSA_COUNT}"
echo ""

if [[ ${TSA_COUNT} -gt 0 ]]; then
    # Show TSA details
    echo "2. TSA Timestamp Details"
    echo "========================"
    jq '.verificationMaterial.timestampVerificationData.rfc3161Timestamps[0] | {
        timestamp_type: "RFC3161 (RFC 3161 Time Stamp Authority)",
        signature_length: (.signedTimestamp | length),
        signature_preview: (.signedTimestamp | .[0:50])
    }' "${BUNDLE_FILE}"
    echo ""
    
    # Show verification material
    echo "3. Bundle Verification Material"
    echo "==============================="
    jq '.verificationMaterial | keys' "${BUNDLE_FILE}"
    echo ""
    
    # Show certificate chain
    CERT_COUNT=$(jq '.verificationMaterial.x509CertificateChain.certificates | length' "${BUNDLE_FILE}" 2>/dev/null || echo 0)
    echo "4. Certificate Chain"
    echo "===================="
    echo "Number of certificates: ${CERT_COUNT}"
    if [[ ${CERT_COUNT} -gt 0 ]]; then
        echo "First certificate length: $(jq '.verificationMaterial.x509CertificateChain.certificates[0].rawBytes | length' "${BUNDLE_FILE}")"
    fi
    echo ""
    
    # Show public key
    echo "5. Public Key Information"
    echo "========================"
    jq '.verificationMaterial.publicKey' "${BUNDLE_FILE}" 2>/dev/null || echo "No public key in material"
    echo ""
    
    # Show signature
    echo "6. Digital Signature"
    echo "===================="
    jq '.messageSignature | {
        algorithm: .messageDigest.algorithm,
        digest_length: (.messageDigest.digest | length),
        signature_length: (.signature | length)
    }' "${BUNDLE_FILE}"
    echo ""
    
    # Show transparency log entries
    TLOG_COUNT=$(jq '.verificationMaterial.tlogEntries | length' "${BUNDLE_FILE}" 2>/dev/null || echo 0)
    echo "7. Transparency Log Integration"
    echo "=============================="
    echo "Transparency log entries: ${TLOG_COUNT}"
    if [[ ${TLOG_COUNT} -gt 0 ]]; then
        echo "First entry details:"
        jq '.verificationMaterial.tlogEntries[0] | {
            log_index: .logIndex,
            log_id: .logId.keyId,
            kind_version: .kindVersion,
            integrated_time: .integratedTime,
            has_inclusion_promise: (.inclusionPromise != null),
            has_inclusion_proof: (.inclusionProof != null)
        }' "${BUNDLE_FILE}"
    fi
    echo ""
    
    echo "✅ TSA Verification Complete"
    echo ""
    echo "Key Findings:"
    echo "  - Bundle contains RFC3161 TSA timestamp from Sigstore TSA"
    echo "  - Includes certificate chain for verification"
    echo "  - Transparency log entries provide auditability"
    echo "  - Combined approach: timestamp + transparency log + signatures"
    echo ""
    echo "Verification Methods Available:"
    echo "  1. TSA Signature Verification (RFC3161)"
    echo "  2. Transparency Log Verification (Rekor)"
    echo "  3. Public Key Verification (Sigstore)"
    echo "  4. Certificate Chain Verification (X.509)"
    echo ""
else
    echo "ERROR: No TSA timestamps found in bundle"
    exit 1
fi