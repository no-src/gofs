#!/usr/bin/env bash

mkdir -p source dest

docker run -it --rm -v "$PWD":/workspace --name running-gofs-local-disk nosrc/gofs:latest \
  gofs -source=./source -dest=./dest
