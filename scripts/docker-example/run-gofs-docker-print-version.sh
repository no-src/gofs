#!/usr/bin/env bash

docker run --rm --name running-gofs-print-version nosrc/gofs:latest \
  gofs -v
