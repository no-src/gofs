#!/usr/bin/env bash

# Remount MinIO using s3fs
# The current directory is integration
cd ../
umount ./integration/minio-data-mount
source ./scripts/minio/mount-minio.sh
