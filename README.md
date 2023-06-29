# Witness Run-Action

# Witness Run GitHub Action

This GitHub Action allows you to create an attestation for your CI process using the Witness tool. It supports optional integration with Sigstore for signing and Archivista for attestation storage and distibution.

## Usage

To use this action, include it in your GitHub workflow YAML file.

### Example

```yaml
permissions:
  id-token: write # This is required for requesting the JWT
  contents: read  # This is required for actions/checkout

name: Example Workflow
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v2

      - name: Witness Run
        uses: testifysec/witness-run-action@v0.1
        with:
          step: test
          use-archivista: true
          use-sigstore: true
          command: touch hello.txt
```

## Using Sigstore and Archivista Flags
This action supports the use of Sigstore and Archivista for creating attestations. By enabling these options, you create a public record of your attestations, which can be useful for transparency and compliance.

### Sigstore
Sigstore is an open-source platform for securely signing software artifacts. When the use-sigstore flag is set to true, this action will use Sigstore for signing the attestation. This creates a publicly verifiable record of the attestation on the Sigstore public instance, sigstore.dev

### Archivista
Archivista is a server that stores and retrieves attestations. When the enable-archivista flag is set to true, this action will use Archivista for storing and retrieving attestations. By default, the attestations are stored on a public Archivista server, archivista.testifysec.io, making the details publicly accessible.  This server also has no guarantees on data availability or itegrity.

### TimeStamping

By default when using Sigstore, this action utilizes FreeTSA, a free and public Timestamp Authority (TSA) service, to provide trusted timestamping for your attestations. Timestamping is a critical aspect of creating non-repudiable and legally binding attestations. FreeTSA offers a reliable and convenient solution for timestamping without the need for setting up and managing your own TSA. When using this action, the timestamp-servers input is set to FreeTSA's service (https://freetsa.org/) by default, ensuring your attestations are properly timestamped with a trusted and publicly verifiable source.

### Privacy Considerations
If you want to keep the details of your attestations private, you can set up and host your own instances of Archivista and Sigstore. This allows you to manage access control and ensure that only authorized users can view the attestation details.

To use your own instances, set the archivista-server input to the URL of your Archivista server, and the fulcio input to the address of your Sigstore instance. Additionally, you'll need to configure the fulcio-oidc-client-id and fulcio-oidc-issuer inputs to match your Sigstore instance's OIDC configuration.

Please consult the documentation for Archivista and Sigstore on how to set up and host your own instances.


### Inputs

| Name                     | Description                                                                                          | Required | Default                               |
| ------------------------ | ---------------------------------------------------------------------------------------------------- | -------- | ------------------------------------- |
| enable-sigstore             | Use Sigstore for attestation. Sets default values for fulcio, fulcio-oidc-client-id, fulcio-oidc-issuer, and timestamp-servers when true | No       | true |
| enable-archivista        | Use Archivista to store or retrieve attestations                                                     | No       | true                                 | true |
| archivista-server        | URL of the Archivista server to store or retrieve attestations                                      | No       | <https://archivista.testifysec.io>      |
| attestations             | Attestations to record, space-separated                                                              | No       | environment git github                      |
| certificate              | Path to the signing key's certificate                                                                | No       |                                       |
| fulcio                   | Fulcio address to sign with                                                                          | No       |                                       |
| fulcio-oidc-client-id    | OIDC client ID to use for authentication                                                             | No       |                                       |
| fulcio-oidc-issuer       | OIDC issuer to use for authentication                                                                | No       |                                       |
| fulcio-token             | Raw token to use for authentication                                                                  | No       |                                       |
| intermediates            | Intermediates that link trust back to a root of trust in the policy, space-separated                | No       |                                       |
| key                      | Path to the signing key                                                                              | No       |                                       |
| outfile                  | File to which to write signed data. Defaults to stdout                                               | No       |                                       |
| product-exclude-glob     | Pattern to use when recording products. Files that match this pattern will be excluded as subjects on the attestation. | No       |                                       |
| product-include-glob     | Pattern to use when recording products. Files that match this pattern will be included as subjects on the attestation. | No       | *                                     |
| spiffe-socket            | Path to the SPIFFE Workload API socket                                                               | No       |                                       |
| step                     | Name of the step being run                                                                           | Yes      |                                       |
| timestamp-servers        | Timestamp Authority Servers to use when signing envelope, space-separated                           | No       |                                       |
| trace                    | Enable tracing for the command                                                                       | No       | false                                 |
| workingdir               | Directory from which commands will run                                                               | No       |                                       |

