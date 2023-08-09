#!/usr/bin/env bash

if ! sudo id sftp_user &>/dev/null; then
  sudo useradd -m sftp_user
  sudo echo "sftp_user:sftp_pwd" | sudo chpasswd
fi

if [ ! -d /sftp-workspace ]; then
  sudo mkdir /sftp-workspace
  sudo chmod 777 /sftp-workspace
fi
