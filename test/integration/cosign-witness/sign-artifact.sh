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

# sign-artifact.sh - Create and sign a test artifact with cosign

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
KEYS_DIR="${SCRIPT_DIR}/keys"
ARTIFACTS_DIR="${SCRIPT_DIR}/artifacts"

echo "=== Step 2: Create and Sign Test Artifact ==="

# Create artifacts directory
mkdir -p "${ARTIFACTS_DIR}"

# Verify keys exist
if [[ ! -f "${KEYS_DIR}/cosign.key" ]]; then
    echo "ERROR: cosign.key not found. Run generate-keys.sh first."
    exit 1
fi

# Create test artifact
TEST_ARTIFACT="${ARTIFACTS_DIR}/test-artifact.txt"
echo "This is a test artifact for Archivista integration" > "${TEST_ARTIFACT}"
echo "Generated at: $(date)" >> "${TEST_ARTIFACT}"
echo "Content hash: $(sha256sum "${TEST_ARTIFACT}" | awk '{print $1}')" >> "${TEST_ARTIFACT}"

echo "Created test artifact: ${TEST_ARTIFACT}"
echo "Content:"
cat "${TEST_ARTIFACT}"
echo ""

# Sign the artifact with cosign to create a Sigstore bundle
BUNDLE_FILE="${ARTIFACTS_DIR}/test-bundle.json"

echo "Signing artifact with cosign..."
# Use --tlog-upload=false to skip transparency log (for local testing)
COSIGN_PASSWORD="" cosign sign-blob \
    --key="${KEYS_DIR}/cosign.key" \
    --bundle="${BUNDLE_FILE}" \
    --tlog-upload=false \
    "${TEST_ARTIFACT}"

echo ""
echo "Bundle created successfully: ${BUNDLE_FILE}"
echo "Bundle size: $(wc -c < "${BUNDLE_FILE}") bytes"
echo ""
echo "Bundle structure:"
jq '.' "${BUNDLE_FILE}" || cat "${BUNDLE_FILE}"
echo ""