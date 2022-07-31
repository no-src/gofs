#!/usr/bin/env bash

mkdir -p sftp-push-client/source sftp-push-client/dest

# depending on your situation, set the SFTP server address
export SFTP_SERVER_ADDR=10.0.4.8
export SFTP_SERVER_REMOTE_PATH=/gofs_sftp_server

export WORKDIR=/workspace

docker run -it --rm -v "$PWD":"$WORKDIR" --name running-gofs-sftp-push-client nosrc/gofs:latest \
  gofs -source="$WORKDIR/sftp-push-client/source" -dest="sftp://$SFTP_SERVER_ADDR:22?local_sync_disabled=false&path=$WORKDIR/sftp-push-client/dest&remote_path=$SFTP_SERVER_REMOTE_PATH" -users="sftp_user|sftp_pwd"
