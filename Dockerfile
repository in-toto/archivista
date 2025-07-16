# Copyright 2022 The Archivista Contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.24.5-alpine@sha256:ddf52008bce1be455fe2b22d780b6693259aaf97b16383b6372f4b22dd33ad66 AS build
WORKDIR /src
RUN apk update && apk add --no-cache file git curl
RUN curl -sSf https://atlasgo.sh | sh
ENV GOMODCACHE /root/.cache/gocache
RUN --mount=target=. --mount=target=/root/.cache,type=cache \
    CGO_ENABLED=0 go build -o /out/archivista -ldflags '-s -d -w' ./cmd/archivista; \
    file /out/archivista | grep "statically linked"

FROM alpine:3.22.1@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1
COPY --from=build /out/archivista /bin/archivista
COPY --from=build /usr/local/bin/atlas /bin/atlas
ADD entrypoint.sh /bin/entrypoint.sh
ADD ent/migrate/migrations /archivista/migrations
RUN mkdir /tmp/archivista
ENTRYPOINT ["sh", "/bin/entrypoint.sh"]
