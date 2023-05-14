#!/bin/bash

set -e

echo "Creating Hydra secret"
kubectl create namespace hydra || true
kubectl apply -f - <<EOF
    apiVersion: v1
    data:
        dsn: bXlzcWw6Ly9yb290OnJvb3RAdGNwKG15c3FsLm15c3FsLnN2Yy5jbHVzdGVyLmxvY2FsOjMzMDYpL2h5ZHJh
        secretsCookie: ZWpUaU1CNUFnek5ZOVFSVk9iWEpOWUJ5M1Ricmd4bjI=
        secretsSystem: QzNBRnZIWlpBM1ZHdmFZZnFxUmFIcjVRdHRlaVlKcVY=
    kind: Secret
    metadata:
        annotations:
            helm.sh/hook: pre-install, pre-upgrade
            helm.sh/hook-delete-policy: before-hook-creation
            helm.sh/hook-weight: "0"
            helm.sh/resource-policy: keep
        labels:
            app.kubernetes.io/instance: hydra
            app.kubernetes.io/managed-by: Helm
            app.kubernetes.io/name: hydra
            app.kubernetes.io/version: v1.11.8
            helm.sh/chart: hydra-0.25.6
        name: hydra
        namespace: hydra
    type: Opaque
EOF
