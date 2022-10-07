#!/usr/bin/env bash

mkdir -p source/encrypt dest/encrypt decrypt_out

export WORKDIR=/workspace

docker run -it --rm -v "$PWD":"$WORKDIR" --name running-gofs-local-disk-with-decrypt nosrc/gofs:latest \
  gofs -source="$WORKDIR/source" -dest="$WORKDIR/dest" -decrypt -decrypt_path="$WORKDIR/dest/encrypt" -decrypt_secret=helloworld -decrypt_out="$WORKDIR/decrypt_out"
