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

# validate-setup.sh - Validate that all prerequisites are met

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ARCHIVISTA_URL="${ARCHIVISTA_URL:-http://localhost:8082}"

echo "Validating test prerequisites..."
echo ""

ERRORS=0

# Check for required binaries
echo "Checking required tools..."

if command -v cosign &> /dev/null; then
    COSIGN_VERSION=$(cosign version 2>&1 | head -1 || echo "unknown")
    echo "  ✓ cosign found: ${COSIGN_VERSION}"
else
    echo "  ✗ cosign not found (install with: brew install cosign or go install github.com/sigstore/cosign/v2/cmd/cosign@latest)"
    ERRORS=$((ERRORS + 1))
fi

if command -v witness &> /dev/null; then
    WITNESS_VERSION=$(witness version 2>&1 | grep -i version || echo "unknown")
    echo "  ✓ witness found: ${WITNESS_VERSION}"
else
    echo "  ✗ witness not found (install with: brew install testifysec/tap/witness or go install github.com/in-toto/go-witness/cmd/witness@latest)"
    ERRORS=$((ERRORS + 1))
fi

ARCHIVISTACTL="/Users/nkennedy/proj/archivista/archivistactl"
if [[ -x "${ARCHIVISTACTL}" ]]; then
    echo "  ✓ archivistactl found at ${ARCHIVISTACTL}"
else
    echo "  ✗ archivistactl not found or not executable at ${ARCHIVISTACTL}"
    echo "    Build with: make archivistactl"
    ERRORS=$((ERRORS + 1))
fi

if command -v jq &> /dev/null; then
    echo "  ✓ jq found"
else
    echo "  ✗ jq not found (install with: brew install jq)"
    ERRORS=$((ERRORS + 1))
fi

if command -v curl &> /dev/null; then
    echo "  ✓ curl found"
else
    echo "  ✗ curl not found"
    ERRORS=$((ERRORS + 1))
fi

echo ""

# Check Archivista availability via GraphQL endpoint
echo "Checking Archivista server..."
if curl -s -f -X POST "${ARCHIVISTA_URL}/query" -H "Content-Type: application/json" -d '{"query":"{__typename}"}' > /dev/null 2>&1; then
    echo "  ✓ Archivista is running at ${ARCHIVISTA_URL}"
else
    echo "  ✗ Archivista is not responding at ${ARCHIVISTA_URL}"
    echo "    Start Archivista with: make run-archivista"
    echo "    Or set ARCHIVISTA_URL environment variable to the correct address"
    ERRORS=$((ERRORS + 1))
fi

echo ""

# Check directory structure
echo "Checking test directory structure..."
if [[ -d "${SCRIPT_DIR}" ]]; then
    echo "  ✓ Test directory exists: ${SCRIPT_DIR}"
else
    echo "  ✗ Test directory not found"
    ERRORS=$((ERRORS + 1))
fi

echo ""

# Summary
if [[ ${ERRORS} -eq 0 ]]; then
    echo "=========================================="
    echo "✓ ALL PREREQUISITES MET"
    echo "=========================================="
    echo ""
    echo "Ready to run integration test!"
    echo "Execute: ./run-test.sh"
    echo ""
    exit 0
else
    echo "=========================================="
    echo "✗ ${ERRORS} PREREQUISITE(S) MISSING"
    echo "=========================================="
    echo ""
    echo "Please install missing tools and ensure Archivista is running."
    echo ""
    exit 1
fi