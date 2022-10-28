# archivist

Archivist helps organizations discover attestations and provenance of their software artifacts.

Archivist is an attestation store with first-class support for Witness attestations but supports any [in-toto](https://in-toto.io) attestation making it work well with other open-source tools that generate in-toto attestations.

## building

```sh
$ docker-compose up --build

$ archivistctl attestation.json
```

## shutting down

```sh
$ docker-compose down
```

## Running archivist out of docker-compose

This application is configured through the environment. The following environment variables can be used:

```sh
KEY                        TYPE             DEFAULT                     REQUIRED    DESCRIPTION
ARCHIVIST_ENABLE_SPIFFE    True or False    TRUE                                    Enable SPIFFE support
ARCHIVIST_LISTEN_ON        URL              unix:///listen.on.socket                url to listen on
ARCHIVIST_LOG_LEVEL        String           INFO                                    Log level
```

Running in a test environment:

```sh
$ go install ./cmd/archivist
$ ARCHIVIST_ENABLE_SPIFFE=false ARCHIVIST_LISTEN_ON=tcp://127.0.0.1:8080 archivist
```

`archivectl` is used to upload and download DSSE objects from the command line. As of now, it only uploads then
downloads the same object to test end to end functionality. This command will be built up in time.

```sh
$ archivistctl file-to-upload-and-downlaod
```
