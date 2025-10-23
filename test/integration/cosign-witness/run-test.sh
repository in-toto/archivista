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

# run-test.sh - Main orchestration script for cosign-witness-archivista integration test

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================="
echo "Cosign-Witness-Archivista Integration Test"
echo "=========================================="
echo ""
echo "This test demonstrates:"
echo "  1. Signing artifacts with cosign"
echo "  2. Storing bundles in Archivista"
echo "  3. Retrieving bundles from Archivista"
echo "  4. Verifying attestations with Witness policies"
echo "  5. Verifying RFC3161 TSA requirements from Witness policy"
echo ""
echo "Starting test at: $(date)"
echo ""

# Set color output for better readability
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

run_step() {
    local step_num=$1
    local step_script=$2
    local step_name=$3

    echo ""
    echo "=========================================="
    echo "STEP ${step_num}: ${step_name}"
    echo "=========================================="
    echo ""

    if bash "${SCRIPT_DIR}/${step_script}"; then
        echo ""
        echo -e "${GREEN}✓ STEP ${step_num} PASSED${NC}"
        return 0
    else
        echo ""
        echo -e "${RED}✗ STEP ${step_num} FAILED${NC}"
        return 1
    fi
}

# Clean up previous test artifacts (optional)
CLEAN="${CLEAN:-false}"
if [[ "${CLEAN}" == "true" ]]; then
    echo "Cleaning previous test artifacts..."
    rm -rf "${SCRIPT_DIR}/artifacts" "${SCRIPT_DIR}/keys"
    echo "Clean complete."
    echo ""
fi

# Run all test steps
FAILED=0

run_step 1 "generate-keys.sh" "Generate Cosign Key Pair" || FAILED=1
run_step 2 "sign-artifact.sh" "Create and Sign Test Artifact" || FAILED=1
run_step 3 "upload-to-archivista.sh" "Upload Bundle to Archivista" || FAILED=1
run_step 4 "download-from-archivista.sh" "Download Bundle from Archivista" || FAILED=1
run_step 5 "create-witness-policy.sh" "Create Witness Policy" || FAILED=1
run_step 6 "verify-with-witness.sh" "Verify with Witness" || FAILED=1
run_step 7 "verify-tsa-policy.sh" "Verify TSA Policy Requirements" || FAILED=1

echo ""
echo "=========================================="
if [[ ${FAILED} -eq 0 ]]; then
    echo -e "${GREEN}ALL TESTS PASSED!${NC}"
    echo "=========================================="
    echo ""
    echo "Summary:"
    echo "  ✓ Cosign keys generated"
    echo "  ✓ Artifact signed with cosign"
    echo "  ✓ Bundle uploaded to Archivista"
    echo "  ✓ Bundle downloaded from Archivista"
    echo "  ✓ Witness policy created"
    echo "  ✓ Verification completed"
    echo "  ✓ TSA policy requirements verified"
    echo ""
    echo "Key findings:"
    echo "  - Archivista successfully stores and retrieves cosign bundles"
    echo "  - Gitoid-based retrieval works correctly"
    echo "  - Downloaded bundles maintain integrity"
    echo "  - Cosign verification works on retrieved bundles"
    echo "  - RFC3161 TSA timestamps present in bundles"
    echo "  - Witness policy correctly configured for TSA validation"
    echo "  - Transparency log integration confirmed"
    echo ""
    echo "Artifacts saved in: ${SCRIPT_DIR}/artifacts/"
    echo "Keys saved in: ${SCRIPT_DIR}/keys/"
    echo ""
    exit 0
else
    echo -e "${RED}SOME TESTS FAILED!${NC}"
    echo "=========================================="
    echo ""
    echo "Please review the output above for details."
    echo ""
    exit 1
fi