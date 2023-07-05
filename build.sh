#!/bin/sh

set -e

DIR=$(dirname "$0")

rm -rf "$DIR/plugin/rootfs"

if docker plugin ls | grep -q oci-secrets; then
  docker plugin rm oci-secrets
fi

docker build -t oci-secrets-rootfs "$DIR"
id=$(docker create oci-secrets-rootfs true)
mkdir -p "$DIR/plugin/rootfs"
docker export "$id" | tar -x -C "$DIR/plugin/rootfs"
docker rm -vf "$id"
docker rmi oci-secrets-rootfs

docker plugin create oci-secrets "$DIR/plugin"
