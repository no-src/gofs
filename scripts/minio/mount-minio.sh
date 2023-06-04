#!/usr/bin/env bash

# Mount MinIO using s3fs
if [ ! -d ./integration/minio-data-mount ]; then
  mkdir -p integration/minio-data-mount
fi
s3fs minio-bucket ./integration/minio-data-mount -o passwd_file=~/.passwd-s3fs -o url=http://127.0.0.1:9000 -o use_path_request_style
