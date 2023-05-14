#!/bin/bash

set -e

minikube ssh -- mkdir -p /home/docker/minio-data /home/docker/mysql-data

## MinIO PVC

kubectl create namespace minio || true

kubectl apply -f - <<EOF
    apiVersion: v1
    kind: PersistentVolume
    metadata:
        namespace: minio
        name: minio-pv-volume
        labels:
            type: local
    spec:
        storageClassName: manual
        capacity:
            storage: 10Gi
        accessModes:
            - ReadWriteOnce
        hostPath:
            path: "/home/docker/minio-data"
EOF

kubectl apply -f - <<EOF
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
        namespace: minio
        name: minio-pvc
    spec:
        storageClassName: manual
        accessModes:
            - ReadWriteOnce
        resources:
            requests:
                storage: 10Gi
EOF

## MySQL PVC

kubectl create namespace mysql || true

kubectl apply -f - <<EOF
    apiVersion: v1
    kind: PersistentVolume
    metadata:
        namespace: mysql
        name: mysql-pv-volume
        labels:
            type: local
    spec:
        storageClassName: manual
        capacity:
            storage: 10Gi
        accessModes:
            - ReadWriteOnce
        hostPath:
            path: "/home/docker/mysql-data"
EOF

kubectl apply -f - <<EOF
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
        namespace: mysql
        name: mysql-pv-claim
    spec:
        storageClassName: manual
        accessModes:
            - ReadWriteOnce
        resources:
            requests:
                storage: 3Gi
EOF
