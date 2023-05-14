#!/bin/bash

set -e

echo "Creating Kratos cofigmap"
kubectl create namespace kratos || true
kubectl apply -f ./k8s/kratos/config-cm.yaml
kubectl apply -f ./k8s/kratos/service-account.yaml