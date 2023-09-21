#!/usr/bin/env bash
# Copyright 2023 The Archivista Contributors
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

set -e

DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

checkprograms() {
  local result=0
  for prog in "$@"
  do
    if ! command -v $prog > /dev/null; then
      printf "$prog is required to run this script. please ensure it is installed and in your PATH\n"
      result=1
    fi
  done

  return $result
}

runtests() {
  go run $DIR/../cmd/archivistactl/main.go store $DIR/*.attestation.json
}

waitForArchivista() {
  echo "Waiting for archivista to be ready..."
  for attempt in $(seq 1 6); do
    sleep 10
    local archivistastate=$(docker compose -f "$DIR/../compose.yml" ps archivista --format json | jq -r '.State')
    if [ "$archivistastate" == "running" ]; then
      break
    fi

    if [[ attempt -eq 6 ]]; then
      echo "timed out waiting for archivista"
      exit 1
    fi
  done
}

if ! checkprograms docker jq ; then
  exit 1
fi

echo "Test mysql..."
docker compose -f "$DIR/../compose.yml" up --build -d
waitForArchivista
runtests
docker compose -f "$DIR/../compose.yml" down -v

echo "Test psql..."
docker compose -f "$DIR/../compose.yml" -f "$DIR/../compose.psql.yml" up --build -d
waitForArchivista
runtests
docker compose -f "$DIR/../compose.yml" -f "$DIR/../compose.psql.yml" down -v