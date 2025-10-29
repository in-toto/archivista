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

# create-tsa-policy.sh - Create Witness policy with TSA verification

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ARTIFACTS_DIR="${SCRIPT_DIR}/artifacts"
KEYS_DIR="${SCRIPT_DIR}/keys"

echo "=== Creating Witness Policy with TSA Verification ==="

BUNDLE_FILE="${ARTIFACTS_DIR}/test-bundle.json"
PUBLIC_KEY_FILE="${KEYS_DIR}/cosign.pub"
POLICY_FILE="${SCRIPT_DIR}/witness-policy-tsa.json"

if [[ ! -f "${BUNDLE_FILE}" ]]; then
    echo "ERROR: Bundle file not found at ${BUNDLE_FILE}"
    exit 1
fi

if [[ ! -f "${PUBLIC_KEY_FILE}" ]]; then
    echo "ERROR: Public key not found at ${PUBLIC_KEY_FILE}"
    exit 1
fi

# Load and encode public key
PUBLIC_KEY=$(cat "${PUBLIC_KEY_FILE}" | base64 -w0)
KEY_ID=$(sha256sum "${PUBLIC_KEY_FILE}" | awk '{print $1}' | base64 -w0)

echo "Creating policy with TSA verification support..."
echo ""

# Create policy with TSA certificate verification
cat > "${POLICY_FILE}" << POLICY
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
  },
  "tsa": {
    "enabled": true,
    "verification": "RFC3161",
    "provider": "sigstore",
    "description": "RFC3161 Time Stamp Authority verification enabled for timestamp integrity"
  }
}
POLICY

echo "✅ TSA-enabled policy created: ${POLICY_FILE}"
echo ""
echo "Policy configuration:"
jq '{tsa: .tsa, root_count: (.roots | length), step_count: (.steps | length)}' "${POLICY_FILE}"
echo ""
echo "TSA Configuration:"
jq '.tsa' "${POLICY_FILE}"
echo ""