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

# download-from-archivista.sh - Download bundle from Archivista using gitoid

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ARTIFACTS_DIR="${SCRIPT_DIR}/artifacts"
ARCHIVISTACTL="/Users/nkennedy/proj/archivista/archivistactl"
ARCHIVISTA_URL="${ARCHIVISTA_URL:-http://localhost:8082}"

echo "=== Step 4: Download Bundle from Archivista ==="

GITOID_FILE="${ARTIFACTS_DIR}/gitoid.txt"

# Verify gitoid exists
if [[ ! -f "${GITOID_FILE}" ]]; then
    echo "ERROR: gitoid.txt not found at ${GITOID_FILE}"
    echo "Run upload-to-archivista.sh first."
    exit 1
fi

GITOID=$(cat "${GITOID_FILE}")
echo "Using gitoid: ${GITOID}"
echo ""

# Verify Archivista is running by testing GraphQL endpoint
echo "Checking Archivista availability at ${ARCHIVISTA_URL}..."
if ! curl -s -f -X POST "${ARCHIVISTA_URL}/query" -H "Content-Type: application/json" -d '{"query":"{__typename}"}' > /dev/null 2>&1; then
    echo "ERROR: Archivista is not responding at ${ARCHIVISTA_URL}"
    exit 1
fi
echo "Archivista is available."
echo ""

# Download bundle using curl to /download endpoint
DOWNLOADED_BUNDLE="${ARTIFACTS_DIR}/downloaded-bundle.json"

echo "Downloading bundle from Archivista..."
curl -s -f "${ARCHIVISTA_URL}/download/${GITOID}" > "${DOWNLOADED_BUNDLE}" 2>&1

if [[ $? -ne 0 ]]; then
    echo "ERROR: Failed to download bundle"
    exit 1
fi

echo "Download complete!"
echo "Downloaded bundle saved to: ${DOWNLOADED_BUNDLE}"
echo ""

# Verify the downloaded bundle
if [[ ! -f "${DOWNLOADED_BUNDLE}" ]]; then
    echo "ERROR: Downloaded bundle not found"
    exit 1
fi

echo "Downloaded bundle size: $(wc -c < "${DOWNLOADED_BUNDLE}") bytes"
echo ""
echo "Downloaded bundle structure:"
jq '.' "${DOWNLOADED_BUNDLE}" || cat "${DOWNLOADED_BUNDLE}"
echo ""

# Compare original and downloaded bundles
ORIGINAL_BUNDLE="${ARTIFACTS_DIR}/test-bundle.json"
echo "Comparing original and downloaded bundles..."

if cmp -s "${ORIGINAL_BUNDLE}" "${DOWNLOADED_BUNDLE}"; then
    echo "SUCCESS: Downloaded bundle matches original bundle exactly!"
else
    echo "Note: Bundles differ (this may be expected due to formatting)"
    echo "Comparing key fields..."

    # Compare signatures
    ORIG_SIG=$(jq -r '.base64Signature // .signature' "${ORIGINAL_BUNDLE}" 2>/dev/null || echo "none")
    DOWN_SIG=$(jq -r '.base64Signature // .signature' "${DOWNLOADED_BUNDLE}" 2>/dev/null || echo "none")

    if [[ "${ORIG_SIG}" == "${DOWN_SIG}" ]]; then
        echo "  Signatures match!"
    else
        echo "  WARNING: Signatures differ"
    fi
fi
echo ""