# Archivista CLI – Quick Usage Guide
This guide explains how to use `archivistactl` to store, search, and retrieve in-toto attestations from an Archivista instance.

1. Store an in-toto attestation in Archivista
```
# Stores an attestation in the Archivista server
$ archivistactl store <attestation.json>

# Example
$ archivistactl store build.attestation.json

# Output
build.attestation.json stored with gitoid 4462a729251af54c7699dbca2f7d5bf5759a5fc6273b3cd606da29d531387c86
```

2. Searching for Attestations
The search command accepts one argument, the subject digest:
```
// Searches the archivista instance for an envelope with a specified subject digest.
// Optionally a collection name can be provided to further constrain results.

$ archivistactl search <algo:digest>

# Example
$ archivistactl search sha256:423da4cff198bbffbe3220ed9510d32ba96698e4b1f654552521d1f541abb6dc

# Output
Gitoid: 4462a729251af54c7699dbca2f7d5bf5759a5fc6273b3cd606da29d531387c86
Collection name: build
Attestations: https://witness.dev/attestations/git/v0.1, https://witness.dev/attestations/environment/v0.1, https://witness.dev/attestations/command-run/v0.1, https://witness.dev/attestations/product/v0.1, https://witness.dev/attestations/material/v0.1

```
3. Retrieve attestations from Archivista
* Retrieve Subjects
```
# Retrieves all subjects on an in-toto statement by the envelope gitoid

$ archivistactl retrieve subjects <gitpoid>

# Example
$ archivistactl retrieve subjects 4462a729251af54c7699dbca2f7d5bf5759a5fc6273b3cd606da29d531387c86

# Output
Name: https://witness.dev/attestations/git/v0.1/committeremail:mswift@mswift.dev
Digests: sha256:408404e7a66b471e5630e801c93af66fb9cb01771982ae90b6f755e104281887
Name: https://witness.dev/attestations/product/v0.1/file:testapp
Digests: gitoid:sha256:gitoid:blob:sha256:473a0f4c3be8a93681a267e3b1e9a7dcda1185436fe141f7749120a303721813, gitoid:sha1:gitoid:blob:sha1:85e3a023c97c8aadace2d8c959535abffbf4e175, sha256:423da4cff198bbffbe3220ed9510d32ba96698e4b1f654552521d1f541abb6dc
Name: https://witness.dev/attestations/git/v0.1/parenthash:aa35c1f4b1d41c87e139c2d333f09117fd0daf4f
Digests: sha256:0bc136f5509e96fc8aa290f175428d643a0e65d8e6b61586ad60e9ec983a3370
Name: https://witness.dev/attestations/git/v0.1/commithash:be20100af602c780deeef50c54f5338662ce917c
Digests: sha1:be20100af602c780deeef50c54f5338662ce917c
Name: https://witness.dev/attestations/git/v0.1/authoremail:snyk-bot@snyk.io
Digests: sha256:ee48369be6072c1a49ba519b2eef9272235b0d925a6e7a338f7ffc12a2ca538e
```

* Retrieve Envelope (DSSE)

```
# Retrieves a dsse envelope by it's gitoid from archivista

$ archivistactl retrieve envelope <gitpoid>
$ archivistactl retrieve envelope 4462a729251af54c7699dbca2f7d5bf5759a5fc6273b3cd606da29d531387c86

# Output
{"payload":"eyJfdHlwZSI6Imh0dHBzOi8vaW4tdG90by5pby9TdGF0ZW1lbnQvdjAuMSIsInN1YmplY3QiOlt7Im5hbWUiOiJodHRwczovL3dpdG5lc3MuZGV2L2F0...
```

---
# Archivista HTTP API Endpoints

1. Get `/v1/download/{gitpoid}`
```
curl <archivista_domain>/v1/download/{gitpoid}

# Example
curl localhost:8082/v1/download/4462a729251af54c7699dbca2f7d5bf5759a5fc6273b3cd606da29d531387c86

# Output
{"payload":"eyJfdHlwZSI6Imh0dHBzOi8vaW4tdG90by5pby9Td...
```

2. Post `/v1/query/`
```
// Query subjects directly
curl -X POST http://localhost:8082/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ subjects { edges { node { name subjectDigests { algorithm value } } } } }"}'
```
* Output
```
{"data":{"subjects":{"edges":[{"node":{"name":"https://witness.dev/attestations/git/v0.1/committeremail:mswift@mswift.dev","subjectDigests":[{"algorithm":"sha256","value":"408404e7a66b471e5630e801c93af66fb9cb01771982ae90b6f755e104281887"}]}},{"node":{"name":"https://witness.dev/attestations/product/v0.1/file:testapp","subjectDigests":[{"algorithm":"gitoid:sha256","value":"gitoid:blob:sha256:473a0f4c3be8a93681a267e3b1e9a7dcda1185436fe141f7749120a303721813"},{"algorithm":"gitoid:sha1","value":"gitoid:blob:sha1:85e3a023c97c8aadace2d8c959535abffbf4e175"},{"algorithm":"sha256","value":"423da4cff198bbffbe3220ed9510d32ba96698e4b1f654552521d1f541abb6dc"}]}},{"node":{"name":"https://witness.dev/attestations/git/v0.1/parenthash:aa35c1f4b1d41c87e139c2d333f09117fd0daf4f","subjectDigests":[{"algorithm":"sha256","value":"0bc136f5509e96fc8aa290f175428d643a0e65d8e6b61586ad60e9ec983a3370"}]}},{"node":{"name":"https://witness.dev/attestations/git/v0.1/commithash:be20100af602c780deeef50c54f5338662ce917c","subjectDigests":[{"algorithm":"sha1","value":"be20100af602c780deeef50c54f5338662ce917c"}]}},{"node":{"name":"https://witness.dev/attestations/git/v0.1/authoremail:snyk-bot@snyk.io","subjectDigests":[{"algorithm":"sha256","value":"ee48369be6072c1a49ba519b2eef9272235b0d925a6e7a338f7ffc12a2ca538e"}]}}]}}}
```
3. Post `/v1/upload`
```
curl -X POST http://localhost:8082/v1/upload \
  -H "Content-Type: application/json" \
  --data-binary "@k8s-att.json"       

# Output
{"gitoid":"72d838472bf801f74dfdc94e4ac8c8c3511da28e0b8af577428114afcf8fcd39"}
```
