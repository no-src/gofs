#!/usr/bin/env bash

docker run -it --rm -v "$PWD":/workspace --name running-gofs-checksum nosrc/gofs:latest \
  gofs -source=/app/gofs -checksum
