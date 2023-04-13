#!/usr/bin/env bash

mkdir -p remote-push-server/source remote-push-server/dest

# depending on your situation, set the gofs server address
export GOFS_SERVER_ADDR=10.0.4.8

export WORKDIR=/workspace

docker run -it --rm -v "$PWD":"$WORKDIR" -w "$WORKDIR" --name go-generate-cert golang:latest \
  go run /usr/local/go/src/crypto/tls/generate_cert.go --host $GOFS_SERVER_ADDR

docker run -it --rm -v "$PWD":"$WORKDIR" -p 8105:8105 -p 443:443 --name running-gofs-remote-push-server nosrc/gofs:latest \
  gofs -source="rs://0.0.0.0:8105?mode=server&local_sync_disabled=true&path=$WORKDIR/remote-push-server/source&fs_server=https://$GOFS_SERVER_ADDR" -dest="$WORKDIR/remote-push-server/dest" -users="gofs|password|rw" -tls_cert_file=cert.pem -tls_key_file=key.pem -push_server -token_secret=mysecret_16bytes
