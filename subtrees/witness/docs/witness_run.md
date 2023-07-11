## witness run

Runs the provided command and records attestations about the execution

```sh
witness run [cmd] [flags]
```

### Options

```sh
      --archivista-server string       URL of the Archivista server to store or retrieve attestations (default "https://archivista.testifysec.io")
  -a, --attestations strings           Attestations to record (default [environment,git])
      --certificate string             Path to the signing key's certificate
      --enable-archivista              Use Archivista to store or retrieve attestations
      --fulcio string                  Fulcio address to sign with
      --fulcio-oidc-client-id string   OIDC client ID to use for authentication
      --fulcio-oidc-issuer string      OIDC issuer to use for authentication
      --fulcio-token string            Raw token to use for authentication
  -h, --help                           help for run
  -i, --intermediates strings          Intermediates that link trust back to a root of trust in the policy
  -k, --key string                     Path to the signing key
  -o, --outfile string                 File to which to write signed data.  Defaults to stdout
      --product-excludeGlob string     Pattern to use when recording products. Files that match this pattern will be excluded as subjects on the attestation. Pass multiple patterns like `{pattern/path/one**,pattern/path/two/**}`.
      --product-includeGlob string     Pattern to use when recording products. Files that match this pattern will be included as subjects on the attestation. Pass multiple patterns like `{pattern/path/one**,pattern/path/two/**}`. (default "*")
      --spiffe-socket string           Path to the SPIFFE Workload API socket
  -s, --step string                    Name of the step being run
      --timestamp-servers strings      Timestamp Authority Servers to use when signing envelope
      --trace                          Enable tracing for the command
  -d, --workingdir string              Directory from which commands will run
```

### Options inherited from parent commands

```sh
  -c, --config string      Path to the witness config file (default ".witness.yaml")
  -l, --log-level string   Level of logging to output (debug, info, warn, error) (default "info")
```

### Notes

- Note that when you use `witness run` during a job that generates a significant number of files, like `npm install`, you may want to consider using `product-excludeGlob` to ignore your node_modules paths when recording products.
- Note that when using `product-excludeGlob`, you may require mulitple paths if you need to ignore multiple dependency locations in a monorepo project.

### SEE ALSO

- [witness](witness.md)  - Collect and verify attestations about your build environments
