# archivist

This application is configured through the environment. The following environment
variables can be used:

```sh
KEY                        TYPE             DEFAULT                     REQUIRED    DESCRIPTION
ARCHIVIST_ENABLE_SPIFFE    True or False    TRUE                                    Enable SPIFFE support
ARCHIVIST_LISTEN_ON        URL              unix:///listen.on.socket                url to listen on
ARCHIVIST_LOG_LEVEL        String           INFO                                    Log level
```

Running in a test environment:
```sh
ARCHIVIST_ENABLE_SPIFFE=false ARCHIVIST_LISTEN_ON=tcp://127.0.0.1:8080 ./archivist
```