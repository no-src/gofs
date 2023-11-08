#!/usr/bin/env bash

export MINIO_ROOT_USER=minio_user
export MINIO_ROOT_PASSWORD=minio_pwd
export MINIO_CONSOLE_ADDRESS=:9001

# Install MinIO (https://github.com/minio/minio)
wget -q https://dl.min.io/server/minio/release/linux-amd64/minio
chmod +x minio

# Run MinIO
mkdir minio-data
setsid ./minio server ./minio-data &

# Install mc (https://github.com/minio/mc)
wget -q https://dl.min.io/client/mc/release/linux-amd64/mc
chmod +x mc

# Set MinIO alias
./mc alias set minio http://127.0.0.1:9000 ${MINIO_ROOT_USER} ${MINIO_ROOT_PASSWORD}

# Create bucket
./mc mb minio/minio-bucket

# Install s3fs (https://github.com/s3fs-fuse/s3fs-fuse)
sudo apt-get install -y s3fs
echo "${MINIO_ROOT_USER}:${MINIO_ROOT_PASSWORD}" >~/.passwd-s3fs
chmod 600 ~/.passwd-s3fs
