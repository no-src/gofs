#!/usr/bin/env bash

mkdir -p remote-disk-client/source remote-disk-client/dest

# depending on your situation, set the gofs server address
export GOFS_SERVER_ADDR=10.0.4.8

export WORKDIR=/workspace

docker run -it --rm -v "$PWD":"$WORKDIR" --name running-gofs-remote-disk-client nosrc/gofs:latest \
  gofs -source="rs://$GOFS_SERVER_ADDR:8105" -dest="$WORKDIR/remote-disk-client/dest" -users="gofs|password" -tls_cert_file=cert.pem
