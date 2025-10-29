# Cosign-Witness-Archivista Integration Test Design

## Overview

This test validates the end-to-end integration of three critical supply chain security tools:
- **Cosign**: Signs artifacts and creates Sigstore bundles
- **Archivista**: Stores and retrieves attestations using content-addressable storage
- **Witness**: Enforces supply chain policies through verification

## Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                         Test Workflow                             │
└──────────────────────────────────────────────────────────────────┘

    1. Key Generation                    2. Artifact Signing
    ┌────────────┐                       ┌─────────────────┐
    │   cosign   │ ──generates──>        │ test-artifact   │
    │generate-key│                       │     .txt        │
    └────────────┘                       └────────┬────────┘
         │                                        │
         v                                        v
    cosign.key                           ┌────────────────┐
    cosign.pub                           │     cosign     │
                                         │   sign-blob    │
                                         └────────┬───────┘
                                                  │
                                                  v
                                          test-bundle.json
                                          (Sigstore format)

    3. Upload to Archivista              4. Download from Archivista
    ┌─────────────────┐                 ┌──────────────────┐
    │ archivistactl   │ ──stores──>     │   Archivista     │
    │     store       │                 │   (via gitoid)   │
    └────────┬────────┘                 └────────┬─────────┘
             │                                    │
             v                                    v
         gitoid: abc123...              downloaded-bundle.json
                                        (Retrieved bundle)

    5. Policy Creation                   6. Verification
    ┌─────────────────┐                 ┌──────────────────┐
    │ Witness Policy  │                 │   cosign verify  │
    │   Generator     │                 │   (baseline)     │
    └────────┬────────┘                 └────────┬─────────┘
             │                                    │
             v                                    v
    witness-policy.json                       ✓ PASS
         │
         v
    ┌─────────────────┐
    │witness verify   │
    │   (policy-      │
    │    based)       │
    └────────┬────────┘
             │
             v
         ✓ PASS
```

## Data Flow

### 1. Signing Phase
```
Input: test-artifact.txt
Key: cosign.key (private)
Process: cosign sign-blob
Output: test-bundle.json (Sigstore bundle)

Bundle Structure:
{
  "base64Signature": "...",
  "cert": "...",
  "rekorBundle": {...},
  "verificationMaterial": {...}
}
```

### 2. Storage Phase
```
Input: test-bundle.json
Process: archivistactl store
Output: gitoid (content-addressable identifier)

Archivista Storage:
- Computes gitoid from bundle content
- Stores bundle in database
- Returns gitoid for retrieval
```

### 3. Retrieval Phase
```
Input: gitoid
Process: archivistactl download
Output: downloaded-bundle.json

Verification:
- Compare with original bundle
- Validate cryptographic integrity
```

### 4. Verification Phase
```
Input: downloaded-bundle.json, cosign.pub, test-artifact.txt
Process:
  a) cosign verify-blob (baseline)
  b) witness verify (policy-based)

Witness Policy:
{
  "publickeys": {"cosign-test-key": {...}},
  "steps": {
    "cosign-sign": {
      "functionaries": [{"type": "publickey", "publickeyid": "cosign-test-key"}]
    }
  }
}
```

## Test Assertions

### Critical Assertions
1. **Key Generation**: Cosign keys are generated successfully
2. **Signing**: Artifact is signed and bundle is created
3. **Upload**: Bundle is stored in Archivista
4. **Gitoid**: Valid gitoid is returned
5. **Download**: Bundle is retrieved using gitoid
6. **Integrity**: Downloaded bundle matches original
7. **Cosign Verification**: Bundle verifies with cosign
8. **Witness Verification**: Bundle verifies with Witness policy

### Integration Points
- **Cosign → Archivista**: Sigstore bundles are stored
- **Archivista → Cosign**: Retrieved bundles verify successfully
- **Gitoid Storage**: Content-addressable retrieval works
- **Witness Policy**: Public key integration works

## File Structure

```
test/integration/cosign-witness/
├── README.md                    # User documentation
├── DESIGN.md                    # This file - technical design
├── Makefile                     # Convenient make targets
│
├── validate-setup.sh            # Prerequisites checker
├── run-test.sh                  # Main orchestrator
├── clean.sh                     # Cleanup script
│
├── generate-keys.sh             # Step 1: Key generation
├── sign-artifact.sh             # Step 2: Signing
├── upload-to-archivista.sh      # Step 3: Upload
├── download-from-archivista.sh  # Step 4: Download
├── create-witness-policy.sh     # Step 5: Policy creation
├── verify-with-witness.sh       # Step 6: Verification
│
├── keys/                        # Generated keys
│   ├── cosign.key              # Private key
│   └── cosign.pub              # Public key
│
├── artifacts/                   # Test artifacts
│   ├── test-artifact.txt       # Signed artifact
│   ├── test-bundle.json        # Original bundle
│   ├── gitoid.txt              # Archivista gitoid
│   ├── downloaded-bundle.json  # Retrieved bundle
│   ├── upload-result.json      # Upload response
│   ├── cosign-verify-result.txt
│   └── witness-verify-result.txt
│
├── witness-policy.json          # Generated policy
└── witness-policy-signed.json   # Signed policy
```

## Key Technical Decisions

### 1. Empty Password for Cosign Keys
**Decision**: Use `COSIGN_PASSWORD=""` for test keys
**Rationale**: Enables automation without interactive prompts
**Production**: Use strong passwords or KMS-backed keys

### 2. Bundle Format
**Decision**: Use Sigstore bundle format from cosign
**Rationale**: Industry standard, Archivista compatible
**Note**: May need format conversion for full Witness integration

### 3. Gitoid-Based Retrieval
**Decision**: Use gitoid as primary identifier
**Rationale**: Content-addressable storage ensures integrity
**Benefit**: Same content always produces same gitoid

### 4. Dual Verification
**Decision**: Verify with both cosign and Witness
**Rationale**: Demonstrates two different verification approaches
**Learning**: Cosign = signature verification, Witness = policy enforcement

### 5. Idempotent Design
**Decision**: Test can run multiple times without cleanup
**Rationale**: Easier development and debugging
**Implementation**: Check for existing artifacts before generation

## Success Criteria

### Minimal Success (Core Integration)
- [ ] Bundle stores successfully in Archivista
- [ ] Gitoid is returned and valid
- [ ] Bundle downloads successfully
- [ ] Downloaded bundle verifies with cosign

### Full Success (Complete Workflow)
- [ ] All above criteria met
- [ ] Witness policy is generated
- [ ] Witness verification succeeds
- [ ] All test steps pass

### Extended Success (Production-Ready)
- [ ] All above criteria met
- [ ] Documentation is comprehensive
- [ ] Error handling is robust
- [ ] Integration points are clearly identified

## Future Enhancements

1. **Format Adapter**: Convert between cosign bundles and Witness attestations
2. **Multi-Artifact**: Test with multiple artifacts and attestations
3. **Policy Evolution**: Test policy updates and versioning
4. **KMS Integration**: Add support for cloud KMS (AWS, GCP, Azure)
5. **Keyless Signing**: Integrate with Fulcio for OIDC-based signing
6. **CI/CD Integration**: Add GitHub Actions / GitLab CI examples
7. **Performance Testing**: Measure upload/download performance
8. **Concurrent Access**: Test multiple simultaneous operations

## Known Limitations

1. **Format Compatibility**: Cosign bundles may need conversion for Witness
2. **Test Keys**: Using weak keys with empty passwords (test only)
3. **Local Storage**: Test assumes local Archivista instance
4. **No Rekor**: Not testing with public transparency log
5. **Single Step**: Only one attestation step (could expand to multi-step)

## Security Considerations

### Test Environment
- Keys are generated locally with empty passwords
- No secrets should be committed to repository
- Test artifacts are ephemeral and local

### Production Deployment
- Use strong passwords or hardware-backed keys
- Integrate with enterprise KMS
- Use Rekor for transparency log
- Implement key rotation policies
- Separate policy signing keys from artifact signing keys

## Debugging

### Enable Verbose Output
```bash
# Add -x to any script for detailed execution
bash -x ./run-test.sh
```

### Inspect Artifacts
```bash
# View bundle structure
jq '.' artifacts/test-bundle.json

# Check gitoid
cat artifacts/gitoid.txt

# Compare bundles
diff <(jq -S '.' artifacts/test-bundle.json) \
     <(jq -S '.' artifacts/downloaded-bundle.json)
```

### Validate Individual Steps
```bash
# Run one step at a time
./generate-keys.sh
./sign-artifact.sh
# ... etc
```

### Check Archivista
```bash
# Health check
curl http://localhost:8082/healthz

# Query by gitoid
GITOID=$(cat artifacts/gitoid.txt)
curl http://localhost:8082/api/v1/attestations/${GITOID}
```

## Testing the Test

This test suite itself should be validated:
1. Run with clean environment
2. Verify all steps pass
3. Manually inspect generated artifacts
4. Re-run without cleanup (test idempotency)
5. Run after Archivista restart (test persistence)
6. Test with invalid data (negative testing)

## References

- [Cosign Documentation](https://docs.sigstore.dev/cosign/overview/)
- [Witness Documentation](https://witness.dev/)
- [Archivista Documentation](https://archivista.dev/)
- [Sigstore Bundle Format](https://github.com/sigstore/protobuf-specs)
- [In-Toto Attestation Framework](https://in-toto.io/)
