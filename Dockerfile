FROM golang:1.18.1-alpine AS build
WORKDIR /src
RUN apk update && apk add --no-cache file git
ENV GOMODCACHE /root/.cache/gocache
RUN --mount=target=. --mount=target=/root/.cache,type=cache \
    CGO_ENABLED=0 go build -o /out/archivist -ldflags '-s -d -w' ./cmd/archivist; \
    file /out/archivist | grep "statically linked"

FROM alpine
COPY --from=build /out/archivist /bin/archivist
RUN mkdir /tmp/archivist
ENTRYPOINT ["/bin/archivist"]
