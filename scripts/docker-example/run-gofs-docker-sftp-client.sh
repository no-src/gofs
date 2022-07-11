#!/usr/bin/env bash

mkdir -p sftp-client/source sftp-client/dest

# depending on your situation, set the SFTP server address
export SFTP_SERVER_ADDR=10.0.4.8
export SFTP_SERVER_REMOTE_PATH=/gofs_sftp_server

export WORKDIR=/workspace

docker run -it --rm -v "$PWD":"$WORKDIR" --name running-gofs-sftp-client nosrc/gofs:latest \
  gofs -source="$WORKDIR/sftp-client/source" -dest="sftp://$SFTP_SERVER_ADDR:22?local_sync_disabled=false&path=$WORKDIR/sftp-client/dest&remote_path=$SFTP_SERVER_REMOTE_PATH" -users="sftp_user|sftp_pwd"
