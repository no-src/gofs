#!/usr/bin/env bash

mkdir -p dest/encrypt decrypt_out

export WORKDIR=/workspace

docker run --rm -v "$PWD":"$WORKDIR" --name running-gofs-local-disk-with-decrypt nosrc/gofs:latest \
  gofs -decrypt -decrypt_path="$WORKDIR/dest/encrypt" -decrypt_secret=mysecret_16bytes -decrypt_out="$WORKDIR/decrypt_out"
