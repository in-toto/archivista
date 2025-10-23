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

# generate-keys.sh - Generate cosign key pair for testing

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
KEYS_DIR="${SCRIPT_DIR}/keys"

echo "=== Step 1: Generate Cosign Key Pair ==="

# Create keys directory
mkdir -p "${KEYS_DIR}"

# Check if keys already exist
if [[ -f "${KEYS_DIR}/cosign.key" && -f "${KEYS_DIR}/cosign.pub" ]]; then
    echo "Keys already exist at ${KEYS_DIR}"
    echo "Public key:"
    cat "${KEYS_DIR}/cosign.pub"
    exit 0
fi

echo "Generating new cosign key pair..."

# Generate key pair (use empty password for automated testing)
# In production, you'd want a real password
COSIGN_PASSWORD="" cosign generate-key-pair --output-key-prefix="${KEYS_DIR}/cosign"

echo ""
echo "Keys generated successfully!"
echo "Private key: ${KEYS_DIR}/cosign.key"
echo "Public key: ${KEYS_DIR}/cosign.pub"
echo ""
echo "Public key content:"
cat "${KEYS_DIR}/cosign.pub"
echo ""