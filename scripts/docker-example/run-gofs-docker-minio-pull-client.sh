#!/usr/bin/env bash

mkdir -p minio-pull-client/source minio-pull-client/dest

# depending on your situation, set the minio server address
export MINIO_SERVER_ADDR=10.0.4.8
export MINIO_BUCKET=minio-bucket

export WORKDIR=/workspace

docker run --rm -v "$PWD":"$WORKDIR" --name running-gofs-minio-pull-client nosrc/gofs:latest \
  gofs -source="minio://$MINIO_SERVER_ADDR:9000?secure=false&remote_path=$MINIO_BUCKET" -dest="$WORKDIR/minio-pull-client/dest" -users="minio_user|minio_pwd" -sync_once
