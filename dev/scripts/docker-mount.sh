#!/bin/sh
#https://jpetazzo.github.io/2015/01/13/docker-mount-dynamic-volumes/

SCRIPT_DIR=$(dirname "$0")
MYSQL_DIR="$SCRIPT_DIR/../.mysql-data"
MINIO_DIR="$SCRIPT_DIR/../.minio-data"

mkdir -p "$MYSQL_DIR"
mkdir -p "$MINIO_DIR"

set -e

CONTAINER=minikube



mount () {
    CONTAINER_PATH="$2"
    HOSTPATH="$1"

    REALPATH=$(readlink --canonicalize $HOSTPATH)
    FILESYS=$(df -P $REALPATH | tail -n 1 | awk '{print $6}')

    while read DEV MOUNT JUNK
    do [ $MOUNT = $FILESYS ] && break
    done </proc/mounts
    [ $MOUNT = $FILESYS ] # Sanity check!

    while read A B C SUBROOT MOUNT JUNK
    do [ $MOUNT = $FILESYS ] && break
    done < /proc/self/mountinfo 
    [ $MOUNT = $FILESYS ] # Moar sanity check!

    SUBPATH=$(echo $REALPATH | sed s,^$FILESYS,,)
    DEVDEC=$(printf "%d %d" $(stat --format "0x%t 0x%T" $DEV))

    docker-enter $CONTAINER sh -c \
	         "[ -b $DEV ] || mknod --mode 0600 $DEV b $DEVDEC"
    docker-enter $CONTAINER mkdir /tmpmnt
    docker-enter $CONTAINER mount $DEV /tmpmnt
    docker-enter $CONTAINER mkdir -p $CONTAINER_PATH
    docker-enter $CONTAINER mount -o bind /tmpmnt/$SUBROOT/$SUBPATH $CONTAINER_PATH
    docker-enter $CONTAINER umount /tmpmnt
    docker-enter $CONTAINER rmdir /tmpmnt
}

mount $MINIO_DIR /home/docker/minio-data
mount $MYSQL_DIR /home/docker/mysql-data



















#REALPATH=$(readlink --canonicalize $HOSTPATH)
# FILESYS=$(df -P $REALPATH | tail -n 1 | awk '{print $6}')

# while read DEV MOUNT JUNK
# do [ $MOUNT = $FILESYS ] && break
# done </proc/mounts
# [ $MOUNT = $FILESYS ] # Sanity check!

# while read A B C SUBROOT MOUNT JUNK
# do [ $MOUNT = $FILESYS ] && break
# done < /proc/self/mountinfo 
# [ $MOUNT = $FILESYS ] # Moar sanity check!

# SUBPATH=$(echo $REALPATH | sed s,^$FILESYS,,)
# DEVDEC=$(printf "%d %d" $(stat --format "0x%t 0x%T" $DEV))

# docker-enter $CONTAINER sh -c \
# 	     "[ -b $DEV ] || mknod --mode 0600 $DEV b $DEVDEC"
# docker-enter $CONTAINER mkdir /tmpmnt
# docker-enter $CONTAINER mount $DEV /tmpmnt
# docker-enter $CONTAINER mkdir -p $CONTPATH
# docker-enter $CONTAINER mount -o bind /tmpmnt/$SUBROOT/$SUBPATH $CONTPATH
# docker-enter $CONTAINER umount /tmpmnt
# docker-enter $CONTAINER rmdir /tmpmnt