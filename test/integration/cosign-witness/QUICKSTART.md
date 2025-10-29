# Quick Start Guide

## One-Liner Test Execution

```bash
# Validate prerequisites and run test
make validate && make test
```

## Prerequisites

1. **Install Tools**:
   ```bash
   brew install cosign jq
   brew install testifysec/tap/witness
   ```

2. **Build archivistactl**:
   ```bash
   cd /Users/nkennedy/proj/archivista
   make archivistactl
   ```

3. **Start Archivista**:
   ```bash
   make run-archivista
   ```

## Running the Test

### Option 1: Using Make
```bash
make validate  # Check prerequisites
make test      # Run full test
make clean     # Clean up artifacts
```

### Option 2: Using Scripts Directly
```bash
./validate-setup.sh  # Check prerequisites
./run-test.sh        # Run full test
./clean.sh           # Clean up artifacts
```

### Option 3: Step-by-Step (for debugging)
```bash
./generate-keys.sh              # Generate cosign keys
./sign-artifact.sh              # Sign test artifact
./upload-to-archivista.sh       # Upload to Archivista
./download-from-archivista.sh   # Download from Archivista
./create-witness-policy.sh      # Create Witness policy
./verify-with-witness.sh        # Verify with Witness
```

## Expected Results

### Validation Output
```
Validating test prerequisites...
Checking required tools...
  ✓ cosign found
  ✓ witness found
  ✓ archivistactl found
  ✓ jq found
  ✓ curl found
Checking Archivista server...
  ✓ Archivista is running at http://localhost:8082
==========================================
✓ ALL PREREQUISITES MET
==========================================
```

### Test Output
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

## Troubleshooting

### Archivista Not Running
```bash
# Check status
curl http://localhost:8082/healthz

# Start Archivista
cd /Users/nkennedy/proj/archivista
make run-archivista
```

### Missing Tools
```bash
# Install cosign
brew install cosign

# Install witness
brew install testifysec/tap/witness

# Install jq
brew install jq

# Build archivistactl
make archivistactl
```

### Test Fails Midway
```bash
# Clean up and start fresh
./clean.sh
./run-test.sh
```

### Permission Errors
```bash
# Make all scripts executable
chmod +x *.sh
```

## Inspecting Results

### View Generated Keys
```bash
cat keys/cosign.pub
```

### View Signed Bundle
```bash
jq '.' artifacts/test-bundle.json
```

### View Gitoid
```bash
cat artifacts/gitoid.txt
```

### Compare Original and Downloaded Bundles
```bash
diff <(jq -S '.' artifacts/test-bundle.json) \
     <(jq -S '.' artifacts/downloaded-bundle.json)
```

### View Witness Policy
```bash
jq '.' witness-policy.json
```

### View Verification Results
```bash
cat artifacts/cosign-verify-result.txt
cat artifacts/witness-verify-result.txt
```

## Environment Variables

```bash
# Use custom Archivista URL
ARCHIVISTA_URL=https://archivista.example.com ./run-test.sh

# Clean before running
CLEAN=true ./run-test.sh
```

## What Gets Created

After running the test, you'll have:

```
keys/
├── cosign.key          # Private signing key
└── cosign.pub          # Public verification key

artifacts/
├── test-artifact.txt              # The artifact that was signed
├── test-bundle.json               # Original cosign bundle
├── gitoid.txt                     # Archivista gitoid
├── downloaded-bundle.json         # Retrieved bundle
├── upload-result.json             # Upload response
├── cosign-verify-result.txt       # Cosign verification output
└── witness-verify-result.txt      # Witness verification output

witness-policy.json                # Generated Witness policy
witness-policy-signed.json         # Signed policy
```

## Next Steps

After running the test successfully:

1. **Inspect Artifacts**: Look at the generated files to understand formats
2. **Read DESIGN.md**: Understand the technical architecture
3. **Modify Policy**: Try changing the Witness policy and re-running verification
4. **Multiple Artifacts**: Extend the test to sign multiple artifacts
5. **CI/CD Integration**: Adapt for your CI/CD pipeline

## Support

For issues or questions:
- Check `README.md` for detailed documentation
- Check `DESIGN.md` for technical details
- Review script output for error messages
- Ensure Archivista is running and accessible
