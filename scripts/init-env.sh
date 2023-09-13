#!/usr/bin/env bash

# https://github.com/quic-go/quic-go/wiki/UDP-Buffer-Sizes
if [ "$(uname -s)" == "Darwin" ]; then
  sudo sysctl -w kern.ipc.maxsockbuf=3014656
else
  sudo sysctl -w net.core.rmem_max=2500000
  sudo sysctl -w net.core.wmem_max=2500000
fi
