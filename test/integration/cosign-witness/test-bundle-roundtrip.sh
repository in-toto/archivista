#!/usr/bin/env bash
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

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "====================================="
echo "Bundle Roundtrip Test"
echo "====================================="

# Check prerequisites
if ! command -v cosign &> /dev/null; then
    echo "Error: cosign not found. Please install cosign."
    exit 1
fi

if ! command -v jq &> /dev/null; then
    echo "Error: jq not found. Please install jq."
    exit 1
fi

# Create test artifact
echo "1. Creating test artifact..."
echo "test data for bundle roundtrip" > artifacts/test-artifact.txt

# Sign with cosign to generate a Sigstore bundle
echo "2. Signing artifact with cosign..."
cosign sign-blob \
    --bundle artifacts/original-bundle.json \
    --key keys/cosign.key \
    --tlog-upload=false \
    --yes \
    artifacts/test-artifact.txt

echo "3. Verifying original bundle structure..."
if ! jq -e '.mediaType | startswith("application/vnd.dev.sigstore.bundle")' artifacts/original-bundle.json > /dev/null; then
    echo "Error: Original bundle is not a valid Sigstore bundle"
    exit 1
fi

# Import bundle to Archivista
echo "4. Importing bundle to Archivista..."
GITOID=$(archivistactl bundle import artifacts/original-bundle.json | grep -oE 'gitoid:[0-9a-f]+')
echo "Imported with $GITOID"

# Get DSSE ID from Archivista
echo "5. Querying for DSSE ID..."
DSSE_ID=$(archivistactl search --gitoid "$GITOID" | jq -r '.[0].dsseIds[0]')
echo "DSSE ID: $DSSE_ID"

# Export bundle from Archivista
echo "6. Exporting bundle from Archivista..."
archivistactl bundle export "$DSSE_ID" -o artifacts/reconstructed-bundle.json

# Verify reconstructed bundle structure
echo "7. Verifying reconstructed bundle structure..."
if ! jq -e '.mediaType | startswith("application/vnd.dev.sigstore.bundle")' artifacts/reconstructed-bundle.json > /dev/null; then
    echo "Error: Reconstructed bundle is not a valid Sigstore bundle"
    exit 1
fi

# Normalize and compare key fields (not byte-for-byte since we reconstruct)
echo "8. Comparing key bundle fields..."

# Extract and compare payload
ORIG_PAYLOAD=$(jq -r '.dsseEnvelope.payload' artifacts/original-bundle.json)
RECON_PAYLOAD=$(jq -r '.dsseEnvelope.payload' artifacts/reconstructed-bundle.json)

if [ "$ORIG_PAYLOAD" != "$RECON_PAYLOAD" ]; then
    echo "Error: Payload mismatch"
    echo "Original payload: $ORIG_PAYLOAD"
    echo "Reconstructed payload: $RECON_PAYLOAD"
    exit 1
fi
echo "✓ Payload matches"

# Extract and compare signature
ORIG_SIG=$(jq -r '.dsseEnvelope.signatures[0].sig' artifacts/original-bundle.json)
RECON_SIG=$(jq -r '.dsseEnvelope.signatures[0].sig' artifacts/reconstructed-bundle.json)

if [ "$ORIG_SIG" != "$RECON_SIG" ]; then
    echo "Error: Signature mismatch"
    exit 1
fi
echo "✓ Signature matches"

# Extract and compare certificate
ORIG_CERT=$(jq -r '.verificationMaterial.certificate.rawBytes // .verificationMaterial.x509CertificateChain.certificates[0].rawBytes' artifacts/original-bundle.json)
RECON_CERT=$(jq -r '.verificationMaterial.certificate.rawBytes // .verificationMaterial.x509CertificateChain.certificates[0].rawBytes' artifacts/reconstructed-bundle.json)

if [ "$ORIG_CERT" != "$RECON_CERT" ]; then
    echo "Error: Certificate mismatch"
    exit 1
fi
echo "✓ Certificate matches"

# Verify reconstructed bundle with cosign
echo "9. Verifying reconstructed bundle with cosign..."
if cosign verify-blob-attestation \
    --bundle artifacts/reconstructed-bundle.json \
    --key keys/cosign.pub \
    --insecure-ignore-tlog \
    artifacts/test-artifact.txt; then
    echo "✓ Reconstructed bundle verifies with cosign"
else
    echo "Error: Reconstructed bundle failed cosign verification"
    exit 1
fi

echo "====================================="
echo "✅ Bundle roundtrip test PASSED"
echo "====================================="
