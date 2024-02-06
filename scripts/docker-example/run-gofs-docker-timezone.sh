#!/usr/bin/env bash

echo "default timezone:"
docker run --rm -e nosrc/gofs:latest date
echo "Asia/Shanghai timezone:"
docker run --rm -e TZ=Asia/Shanghai nosrc/gofs:latest date
