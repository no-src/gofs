#!/usr/bin/env bash

if ! sudo id sftp_user &>/dev/null; then
  sudo useradd -m sftp_user
  sudo echo "sftp_user:sftp_pwd" | sudo chpasswd
fi

if [ ! -d /sftp-workspace ]; then
  sudo mkdir /sftp-workspace
  sudo chmod 777 /sftp-workspace
fi

# init ssh_config
mkdir ~/.ssh
ssh-keygen -t rsa -b 4096 -q -N "" -f ~/.ssh/id_rsa
ln -s $(pwd)/integration/testdata/ssh/ssh_config ~/.ssh/config

sshpass -v -p 'sftp_pwd' ssh-copy-id -o StrictHostKeyChecking=no sftp_user@127.0.0.1

ls -alh ~/.ssh
ls -alh /home/sftp_user/.ssh
