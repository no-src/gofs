#!/usr/bin/env bash

sudo useradd -m sftp_user
sudo echo "sftp_user:sftp_pwd" | sudo chpasswd

if [ ! -d /sftp-workspace ]; then
  sudo mkdir /sftp-workspace
  sudo chmod 777 /sftp-workspace
fi