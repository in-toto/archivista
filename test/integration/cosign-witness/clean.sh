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

# clean.sh - Clean up test artifacts

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Cleaning test artifacts..."

# Remove artifacts directory
if [[ -d "${SCRIPT_DIR}/artifacts" ]]; then
    rm -rf "${SCRIPT_DIR}/artifacts"
    echo "  ✓ Removed artifacts/"
fi

# Remove keys directory
if [[ -d "${SCRIPT_DIR}/keys" ]]; then
    rm -rf "${SCRIPT_DIR}/keys"
    echo "  ✓ Removed keys/"
fi

# Remove generated policy files
if [[ -f "${SCRIPT_DIR}/witness-policy.json" ]]; then
    rm -f "${SCRIPT_DIR}/witness-policy.json"
    echo "  ✓ Removed witness-policy.json"
fi

if [[ -f "${SCRIPT_DIR}/witness-policy-signed.json" ]]; then
    rm -f "${SCRIPT_DIR}/witness-policy-signed.json"
    echo "  ✓ Removed witness-policy-signed.json"
fi

echo ""
echo "Cleanup complete. Run ./run-test.sh to start fresh."