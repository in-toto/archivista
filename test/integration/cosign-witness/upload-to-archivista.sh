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

# upload-to-archivista.sh - Upload cosign bundle to Archivista

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ARTIFACTS_DIR="${SCRIPT_DIR}/artifacts"
ARCHIVISTACTL="/Users/nkennedy/proj/archivista/archivistactl"
ARCHIVISTA_URL="${ARCHIVISTA_URL:-http://localhost:8082}"

echo "=== Step 3: Upload Bundle to Archivista ==="

BUNDLE_FILE="${ARTIFACTS_DIR}/test-bundle.json"

# Verify bundle exists
if [[ ! -f "${BUNDLE_FILE}" ]]; then
    echo "ERROR: Bundle file not found at ${BUNDLE_FILE}"
    echo "Run sign-artifact.sh first."
    exit 1
fi

# Verify Archivista is running by testing GraphQL endpoint
echo "Checking Archivista availability at ${ARCHIVISTA_URL}..."
if ! curl -s -f -X POST "${ARCHIVISTA_URL}/query" -H "Content-Type: application/json" -d '{"query":"{__typename}"}' > /dev/null 2>&1; then
    echo "ERROR: Archivista is not responding at ${ARCHIVISTA_URL}"
    echo "Please ensure Archivista is running."
    exit 1
fi
echo "Archivista is available."
echo ""

# Upload bundle using curl to /upload endpoint
# Note: /upload endpoint expects plain JSON body, not multipart form data
echo "Uploading bundle to Archivista via /upload endpoint..."
UPLOAD_OUTPUT="${ARTIFACTS_DIR}/upload-result.json"

curl -s -X POST "${ARCHIVISTA_URL}/upload" \
    -H "Content-Type: application/json" \
    -d @"${BUNDLE_FILE}" \
    > "${UPLOAD_OUTPUT}" 2>&1

echo "Upload complete!"
echo ""
echo "Upload result:"
cat "${UPLOAD_OUTPUT}"
echo ""

# Extract gitoid from response
GITOID=$(jq -r '.gitoid // .Gitoid // empty' "${UPLOAD_OUTPUT}" 2>/dev/null || grep -oE 'gitoid:[a-f0-9]+' "${UPLOAD_OUTPUT}" | cut -d: -f2 || echo "")

if [[ -z "${GITOID}" ]]; then
    echo "WARNING: Could not extract gitoid from response. Attempting alternative methods..."
    # Try to extract from any field containing hash-like string
    GITOID=$(cat "${UPLOAD_OUTPUT}" | grep -oE '[a-f0-9]{64}' | head -1 || echo "")
fi

if [[ -n "${GITOID}" ]]; then
    echo "Gitoid: ${GITOID}"
    echo "${GITOID}" > "${ARTIFACTS_DIR}/gitoid.txt"
    echo "Gitoid saved to ${ARTIFACTS_DIR}/gitoid.txt"
else
    echo "ERROR: Failed to extract gitoid from upload response"
    exit 1
fi
echo ""