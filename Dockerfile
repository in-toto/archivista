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

FROM golang:1.22.1-alpine@sha256:0466223b8544fb7d4ff04748acc4d75a608234bf4e79563bff208d2060c0dd79 AS build
WORKDIR /src
RUN apk update && apk add --no-cache file git curl
RUN curl -sSf https://atlasgo.sh | sh
ENV GOMODCACHE /root/.cache/gocache
RUN --mount=target=. --mount=target=/root/.cache,type=cache \
    CGO_ENABLED=0 go build -o /out/archivista -ldflags '-s -d -w' ./cmd/archivista; \
    file /out/archivista | grep "statically linked"

FROM alpine:3.19.1@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b
COPY --from=build /out/archivista /bin/archivista
COPY --from=build /usr/local/bin/atlas /bin/atlas
ADD entrypoint.sh /bin/entrypoint.sh
ADD ent/migrate/migrations /archivista/migrations
RUN mkdir /tmp/archivista
ENTRYPOINT ["sh", "/bin/entrypoint.sh"]
