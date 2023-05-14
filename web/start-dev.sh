#!/bin/bash

set -e

working_dir=$(pwd)

script_dir=$(dirname "$0")
cd "$script_dir"
running_container=""


function buildandrun() {
    docker stop web-dev || true

    DOCKER_BUILDKIT=1 docker build -f Dockerfile.dev -t web:dev .
    running_container=$(docker run --rm --name web-dev -d -p 8077:8077 -p 1234:1234 -v "$PWD":/src -v "$PWD"/config.yaml:/etc/web/config.yaml web:dev)
    docker logs -f $running_container &
    #stop on ctrl+c
    trap "echo 'stopping container' && cleanup" SIGINT
    trap "echo 'stopping container' && cleanup" SIGTERM
    trap "echo 'stopping container' && cleanup" EXIT
    wait
}

function watch() {
  while true; do
    inotifywait -e close_write package.json Dockerfile.dev config.yaml start-dev.sh
    echo "Restarting..."
    docker stop "$running_container"
    buildandrun
  done
}

function cleanup() {
  docker stop "$running_container"
  cd "$working_dir"
  exit 0
}


buildandrun
watch
cleanup
