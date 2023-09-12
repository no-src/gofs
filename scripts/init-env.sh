#!/usr/bin/env bash

# https://github.com/quic-go/quic-go/wiki/UDP-Buffer-Sizes
if [ "$(uname -s)" == "Darwin" ]; then
  sysctl -w kern.ipc.maxsockbuf=3014656
else
  sysctl -w net.core.rmem_max=2500000
  sysctl -w net.core.wmem_max=2500000
fi
