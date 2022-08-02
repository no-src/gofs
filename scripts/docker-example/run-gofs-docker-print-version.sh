#!/usr/bin/env bash

docker run -it --rm --name running-gofs-print-version nosrc/gofs:latest \
  gofs -v
