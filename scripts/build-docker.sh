#!/usr/bin/env bash

# usage
# build the latest image
# ./scripts/build-docker.sh
# build the image with a specified tag
# ./scripts/build-docker.sh v0.6.0

# update git repository
# git pull --no-rebase

# update the latest golang image
docker pull golang:alpine

# set GOPROXY environment variable
# GOPROXY=https://goproxy.cn

# set the gofs docker image name by GOFS_IMAGE_NAME variable
GOFS_IMAGE_NAME=nosrc/gofs

# set the gofs docker image tag by GOFS_IMAGE_TAG variable
GOFS_IMAGE_TAG=latest

# reset GOFS_IMAGE_TAG to the value of the first parameter provided by the user
if [ -n "$1" ]; then
  GOFS_IMAGE_TAG=$1
fi

# remove the existing old image
docker rmi -f $GOFS_IMAGE_NAME:$GOFS_IMAGE_TAG

# build Dockerfile
docker build --build-arg GOPROXY=$GOPROXY -t $GOFS_IMAGE_NAME:$GOFS_IMAGE_TAG .

# remove dangling images
docker image prune -f

# run a container to print the gofs version
docker run -it --rm --name running-gofs-version $GOFS_IMAGE_NAME:$GOFS_IMAGE_TAG gofs -v

# push the image to the DockerHub
# docker push $GOFS_IMAGE_NAME:$GOFS_IMAGE_TAG
