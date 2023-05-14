#!/bin/bash

set -e

echo "Creating Kratos secret"
kubectl create namespace kratos || true
kubectl apply -f - <<EOF
    apiVersion: v1
    data:
        dsn: bXlzcWw6Ly9yb290OnJvb3RAdGNwKG15c3FsLm15c3FsLnN2Yy5jbHVzdGVyLmxvY2FsOjMzMDYpL2tyYXRvcw==
        smtpConnectionURI: c210cHM6Ly90ZXN0OnRlc3RAbWFpbGhvZy5tYWlsaG9nLnN2Yy5jbHVzdGVyLmxvY2FsLz9za2lwX3NzbF92ZXJpZnk9dHJ1ZQ==
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
            app.kubernetes.io/instance: kratos
            app.kubernetes.io/managed-by: Helm
            app.kubernetes.io/name: kratos
            app.kubernetes.io/version: v0.10.1
            helm.sh/chart: kratos-0.26.4
        name: kratos
        namespace: kratos
    type: Opaque
EOF
