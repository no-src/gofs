#!/usr/bin/env bash

mkdir -p source/encrypt dest/encrypt decrypt_out

export WORKDIR=/workspace

docker run -it --rm -v "$PWD":"$WORKDIR" --name running-gofs-local-disk-with-encrypt nosrc/gofs:latest \
  gofs -source="$WORKDIR/source" -dest="$WORKDIR/dest" -encrypt -encrypt_path="$WORKDIR/source/encrypt" -encrypt_secret=helloworld
