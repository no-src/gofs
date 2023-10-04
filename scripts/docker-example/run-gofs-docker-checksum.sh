#!/usr/bin/env bash

export WORKDIR=/workspace

docker run --rm -v "$PWD":"$WORKDIR" --name running-gofs-checksum nosrc/gofs:latest \
  gofs -source=/app/gofs -checksum
