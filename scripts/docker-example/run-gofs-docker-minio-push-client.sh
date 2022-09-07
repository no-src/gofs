#!/usr/bin/env bash

mkdir -p minio-push-client/source minio-push-client/dest

# depending on your situation, set the minio server address
export MINIO_SERVER_ADDR=10.0.4.8
export MINIO_BUCKET=minio-bucket

export WORKDIR=/workspace

docker run -it --rm -v "$PWD":"$WORKDIR" --name running-gofs-minio-push-client nosrc/gofs:latest \
  gofs -source="$WORKDIR/minio-push-client/source" -dest="minio://$MINIO_SERVER_ADDR:9000?local_sync_disabled=false&path=$WORKDIR/minio-push-client/dest&remote_path=$MINIO_BUCKET" -users="minio_user|minio_pwd"
