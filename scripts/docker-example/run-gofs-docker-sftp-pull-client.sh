#!/usr/bin/env bash

mkdir -p sftp-pull-client/source sftp-pull-client/dest

# depending on your situation, set the SFTP server address
export SFTP_SERVER_ADDR=10.0.4.8
export SFTP_SERVER_REMOTE_PATH=/gofs_sftp_server

export WORKDIR=/workspace

docker run -it --rm -v "$PWD":"$WORKDIR" --name running-gofs-sftp-pull-client nosrc/gofs:latest \
  gofs -source="sftp://$SFTP_SERVER_ADDR:22?remote_path=$SFTP_SERVER_REMOTE_PATH&ssh_user=sftp_user&ssh_pass=sftp_pwd" -dest="$WORKDIR/sftp-pull-client/dest" -sync_once
