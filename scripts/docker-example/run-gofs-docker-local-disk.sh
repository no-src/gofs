#!/usr/bin/env bash

mkdir -p source dest

export WORKDIR=/workspace

docker run --rm -v "$PWD":"$WORKDIR" --name running-gofs-local-disk nosrc/gofs:latest \
  gofs -source="$WORKDIR/source" -dest="$WORKDIR/dest"
