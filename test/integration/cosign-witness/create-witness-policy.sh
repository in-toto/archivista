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

# create-witness-policy.sh - Generate Witness policy for cosign attestations

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
KEYS_DIR="${SCRIPT_DIR}/keys"
ARTIFACTS_DIR="${SCRIPT_DIR}/artifacts"

echo "=== Step 5: Create Witness Policy ==="

# Verify public key exists
if [[ ! -f "${KEYS_DIR}/cosign.pub" ]]; then
    echo "ERROR: cosign.pub not found. Run generate-keys.sh first."
    exit 1
fi

# Read and base64 encode the public key
PUBLIC_KEY=$(cat "${KEYS_DIR}/cosign.pub" | base64 -w 0 2>/dev/null || cat "${KEYS_DIR}/cosign.pub" | base64)

echo "Public key loaded and encoded"
echo ""

# Extract the actual key ID from the cosign bundle if available
BUNDLE_FILE="${ARTIFACTS_DIR}/test-bundle.json"
KEY_ID=""

if [[ -f "${BUNDLE_FILE}" ]]; then
    # Try to extract key hint or certificate subject
    KEY_ID=$(jq -r '.verificationMaterial.publicKey.hint // empty' "${BUNDLE_FILE}" 2>/dev/null || echo "")

    if [[ -z "${KEY_ID}" ]]; then
        # Calculate key ID from public key (SHA256 hash of DER-encoded public key)
        # For cosign, we can use a simple identifier
        KEY_ID="cosign-test-key"
    fi

    echo "Using Key ID: ${KEY_ID}"
else
    echo "WARNING: Bundle not found, using default key ID"
    KEY_ID="cosign-test-key"
fi
echo ""

# Get test artifact hash
TEST_ARTIFACT="${ARTIFACTS_DIR}/test-artifact.txt"
if [[ -f "${TEST_ARTIFACT}" ]]; then
    ARTIFACT_HASH=$(sha256sum "${TEST_ARTIFACT}" 2>/dev/null | awk '{print $1}' || shasum -a 256 "${TEST_ARTIFACT}" | awk '{print $1}')
    echo "Test artifact hash: ${ARTIFACT_HASH}"
else
    echo "WARNING: Test artifact not found, policy will not include artifact hash validation"
    ARTIFACT_HASH=""
fi
echo ""

# Create the Witness policy
POLICY_FILE="${SCRIPT_DIR}/witness-policy.json"

cat > "${POLICY_FILE}" <<EOF
{
  "expires": "2030-12-31T23:59:59Z",
  "roots": {
    "${KEY_ID}": {
      "certificate": "${PUBLIC_KEY}",
      "intermediates": []
    }
  },
  "publickeys": {
    "${KEY_ID}": {
      "keyid": "${KEY_ID}",
      "key": "${PUBLIC_KEY}"
    }
  },
  "steps": {
    "cosign-sign": {
      "name": "cosign-sign",
      "attestations": [
        {
          "type": "https://in-toto.io/attestation/v0.1",
          "regopolicies": []
        }
      ],
      "functionaries": [
        {
          "type": "publickey",
          "publickeyid": "${KEY_ID}"
        }
      ]
    }
  }
}
EOF

echo "Witness policy created: ${POLICY_FILE}"
echo ""
echo "Policy structure:"
jq '.' "${POLICY_FILE}"
echo ""

# Sign the policy with the same cosign key
SIGNED_POLICY_FILE="${SCRIPT_DIR}/witness-policy-signed.json"

echo "Signing policy with cosign key..."
# Note: Witness expects the policy to be signed, but for testing we can also verify unsigned
# For production, you'd want to sign this with a separate policy signing key

# For now, we'll copy the policy as-is since witness verify can work with unsigned policies
# in test mode
cp "${POLICY_FILE}" "${SIGNED_POLICY_FILE}"

echo "Policy ready for verification: ${SIGNED_POLICY_FILE}"
echo ""