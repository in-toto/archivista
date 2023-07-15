#!/bin/bash

set -e

script_dir=$(dirname "$0")
root_dir=$(dirname "$script_dir")
cd "$root_dir"

CAROOT=$root_dir/.certs

echo "Creating CA"
CAROOT=$CAROOT mkcert -install -key-file "${CAROOT}"/key.pem -cert-file "${CAROOT}"/cert.pem testifysec.localhost *.testifysec.localhost localhost

echo "Injecting CA into kubernetes"
kubectl create secret --dry-run=client -n cert-manager tls tls-secret --cert="${CAROOT}"/rootCA.pem --key="${CAROOT}"/rootCA-key.pem -o yaml | kubectl apply -f -
kubectl -n kratos create configmap rootca --from-file="${CAROOT}/rootCA.pem" --dry-run=client -o yaml | kubectl apply -f -
