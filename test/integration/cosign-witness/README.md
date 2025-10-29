# Cosign-Witness-Archivista Integration Test

This integration test demonstrates the complete workflow of:
1. Generating cosign keys
2. Signing artifacts with cosign (creating Sigstore bundles)
3. Uploading bundles to Archivista
4. Downloading bundles from Archivista
5. Creating Witness policies that accept cosign signatures
6. Verifying cosign-signed attestations with Witness

## Prerequisites

### Required Tools
- **cosign** - Install with: `brew install cosign` or `go install github.com/sigstore/cosign/v2/cmd/cosign@latest`
- **witness** - Install with: `brew install testifysec/tap/witness` or `go install github.com/in-toto/go-witness/cmd/witness@latest`
- **archivistactl** - Build with: `make archivistactl` from project root
- **jq** - Install with: `brew install jq`
- **curl** - Usually pre-installed

### Running Services
- Archivista server must be running (default: http://localhost:8082)
- Start with: `make run-archivista` or your preferred method

### Validate Prerequisites
Run the validation script to check all requirements:
```bash
./validate-setup.sh
```

## Quick Start

```bash
# Validate prerequisites
./validate-setup.sh

# Run the full integration test
./run-test.sh

# Clean up test artifacts (optional)
./clean.sh
```

## What This Proves

This test proves that:
- **Cosign bundles** (Sigstore format) can be stored in Archivista
- **Attestations** signed by cosign can be retrieved from Archivista using gitoid
- **Witness policies** can be created to verify cosign-signed attestations
- **Gitoid-based storage/retrieval** mechanism works correctly
- **Downloaded bundles** maintain cryptographic integrity
- **Cosign verification** works on bundles retrieved from Archivista

## Test Components

### Scripts (in execution order)
1. **validate-setup.sh** - Validate all prerequisites are met
2. **run-test.sh** - Main test orchestration script (runs all steps)
3. **generate-keys.sh** - Generate cosign key pair (cosign.key, cosign.pub)
4. **sign-artifact.sh** - Create test artifact and sign with cosign
5. **upload-to-archivista.sh** - Upload Sigstore bundle to Archivista
6. **download-from-archivista.sh** - Download bundle using gitoid
7. **create-witness-policy.sh** - Generate Witness policy for verification
8. **verify-with-witness.sh** - Verify attestation with Witness policy
9. **clean.sh** - Clean up test artifacts

### Generated Artifacts
- **keys/** - Cosign key pair (cosign.key, cosign.pub)
- **artifacts/** - Test artifacts and bundles
  - test-artifact.txt - The artifact being signed
  - test-bundle.json - Original cosign bundle
  - gitoid.txt - Gitoid returned by Archivista
  - downloaded-bundle.json - Bundle retrieved from Archivista
  - upload-result.json - Upload response from Archivista
  - cosign-verify-result.txt - Cosign verification output
  - witness-verify-result.txt - Witness verification output
- **witness-policy.json** - Generated Witness policy
- **witness-policy-signed.json** - Signed policy (ready for verification)

## Detailed Test Flow

### Step 1: Generate Keys
Generates a cosign key pair for signing. Uses empty password for automated testing.
```bash
./generate-keys.sh
```

### Step 2: Sign Artifact
Creates a test artifact and signs it with cosign, producing a Sigstore bundle.
```bash
./sign-artifact.sh
```

### Step 3: Upload to Archivista
Uploads the Sigstore bundle to Archivista and captures the gitoid.
```bash
./upload-to-archivista.sh
```

### Step 4: Download from Archivista
Retrieves the bundle using the gitoid and compares with original.
```bash
./download-from-archivista.sh
```

### Step 5: Create Witness Policy
Generates a Witness policy that accepts cosign signatures with the test public key.
```bash
./create-witness-policy.sh
```

### Step 6: Verify with Witness
Verifies the downloaded bundle using both cosign and Witness.
```bash
./verify-with-witness.sh
```

## Expected Output

All steps should succeed with output similar to:
```
==========================================
Cosign-Witness-Archivista Integration Test
==========================================

STEP 1: Generate Cosign Key Pair
✓ STEP 1 PASSED

STEP 2: Create and Sign Test Artifact
✓ STEP 2 PASSED

STEP 3: Upload Bundle to Archivista
✓ STEP 3 PASSED

STEP 4: Download Bundle from Archivista
✓ STEP 4 PASSED

STEP 5: Create Witness Policy
✓ STEP 5 PASSED

STEP 6: Verify with Witness
✓ STEP 6 PASSED

==========================================
ALL TESTS PASSED!
==========================================
```

## Integration Points Validated

This test validates the following integration points:

1. **Cosign → Archivista**: Sigstore bundles can be stored in Archivista
2. **Archivista → Cosign**: Retrieved bundles are cryptographically intact
3. **Archivista gitoid**: Content-addressable storage works correctly
4. **Witness Policy**: Policies can reference cosign public keys
5. **End-to-End**: Complete supply chain attestation flow

## Troubleshooting

### Archivista not responding
- Ensure Archivista is running: `curl http://localhost:8082/healthz`
- Check if correct port: Set `ARCHIVISTA_URL` environment variable
- Start Archivista: `make run-archivista`

### Missing tools
- Run `./validate-setup.sh` to identify missing prerequisites
- Install missing tools using the commands provided

### Permission errors
- Ensure all .sh files are executable: `chmod +x *.sh`

### Clean start
- Run `./clean.sh` to remove all generated artifacts
- Re-run with `./run-test.sh`

## Environment Variables

- **ARCHIVISTA_URL** - Archivista server URL (default: http://localhost:8082)
- **CLEAN** - Set to "true" to clean artifacts before running: `CLEAN=true ./run-test.sh`

## Use Cases

This integration test is useful for:
- Validating Archivista installation and configuration
- Testing cosign bundle storage and retrieval
- Demonstrating supply chain security workflows
- CI/CD integration testing
- Debugging attestation format compatibility

## Notes

- Keys are generated with empty passwords for automation - use real passwords in production
- The test demonstrates storage/retrieval but format conversion may be needed for full Witness integration
- All artifacts are stored locally in the test directory for inspection
- The test is idempotent - can be run multiple times
