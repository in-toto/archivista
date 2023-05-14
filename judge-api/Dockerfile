FROM golang:1.19.1-alpine AS build
WORKDIR /src
RUN apk update && apk add --no-cache file git
ENV GOMODCACHE /root/.cache/gocache
RUN --mount=target=. --mount=target=/root/.cache,type=cache \
    CGO_ENABLED=0 go build -o /out/judge-api -ldflags '-s -d -w' ./cmd/server; \
    file /out/judge-api | grep "statically linked"

FROM alpine
COPY --from=build /out/judge-api /bin/judge-api
ENTRYPOINT ["/bin/judge-api"]