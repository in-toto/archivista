# Copyright 2022 The Archivist Contributors
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

FROM golang:1.19.1-alpine AS build
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
