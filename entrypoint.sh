#!/bin/bash
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

ARCHIVISTA_ENABLE_SQL_STORE=$(echo ${ARCHIVISTA_ENABLE_SQL_STORE} | tr '[:lower:]' '[:upper:]')

if [ "${ARCHIVISTA_ENABLE_SQL_STORE}" = "FALSE" ]; then
  echo "Skipping migrations"
else
  if [[ -z $ARCHIVISTA_SQL_STORE_BACKEND ]]; then
    SQL_TYPE="MYSQL"
  else
    SQL_TYPE=$(echo "$ARCHIVISTA_SQL_STORE_BACKEND" | tr '[:lower:]' '[:upper:]')
  fi
  case $SQL_TYPE in
  MYSQL)
    if [[ -z $ARCHIVISTA_SQL_STORE_CONNECTION_STRING ]]; then
      ARCHIVISTA_SQL_STORE_CONNECTION_STRING="root:example@db/testify"
    fi
    echo "Running migrations for MySQL"
    atlas migrate apply --dir "file:///archivista/migrations/mysql" --url "mysql://$ARCHIVISTA_SQL_STORE_CONNECTION_STRING"
    atlas_rc=$?
    ;;
  PSQL)
    echo "Running migrations for Postgres"
    atlas migrate apply --dir "file:///archivista/migrations/pgsql" --url "$ARCHIVISTA_SQL_STORE_CONNECTION_STRING"
    atlas_rc=$?
    ;;
  *)
    echo "Unknown SQL backend: $ARCHIVISTA_SQL_STORE_BACKEND"
    exit 1
    ;;
  esac

  if [[ $atlas_rc -ne 0 ]]; then
    echo "Failed to apply migrations"
    exit 1
  fi
fi

/bin/archivista
