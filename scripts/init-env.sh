#!/usr/bin/env bash

# https://github.com/quic-go/quic-go/wiki/UDP-Buffer-Sizes
sysctl -w net.core.rmem_max=2500000
sysctl -w net.core.wmem_max=2500000