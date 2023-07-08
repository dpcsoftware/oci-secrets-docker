#!/bin/sh

set -e

DIR=$(dirname "$0")

if [ $# -eq 0 ]; then
  PLATFORM="amd64"
else
  PLATFORM=$1
fi

PLUGIN_NAME=dpcsoftware/oci-secrets
VERSION=0.1

rm -rf "$DIR/plugin/rootfs"

if docker plugin ls | grep -q "$PLUGIN_NAME:$VERSION-$PLATFORM"; then
  docker plugin rm "$PLUGIN_NAME"
fi

docker buildx build --platform $PLATFORM -t oci-secrets-rootfs "$DIR"
id=$(docker create oci-secrets-rootfs true)
mkdir -p "$DIR/plugin/rootfs"
docker export "$id" | tar -x -C "$DIR/plugin/rootfs"
docker rm -vf "$id"
docker rmi oci-secrets-rootfs

docker plugin create "$PLUGIN_NAME:$VERSION-$PLATFORM" "$DIR/plugin"
