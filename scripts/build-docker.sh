#!/usr/bin/env bash

# usage
# build the latest image
# ./scripts/build-docker.sh
# build the image with a specified tag
# ./scripts/build-docker.sh v0.8.0

echo "current branch is $(git branch --show-current)"

# update git repository
# git pull --no-rebase

# update the latest golang image
docker pull golang:latest

# set GOPROXY environment variable
# GOPROXY=https://goproxy.cn

SOFT_NAME=gofs
# docker image name
SOFT_IMAGE_NAME=nosrc/${SOFT_NAME}
# docker image tag
SOFT_IMAGE_TAG=latest

# reset SOFT_IMAGE_TAG to the value of the first parameter provided by the user
if [ -n "$1" ]; then
  SOFT_IMAGE_TAG=$1
fi

# remove the existing old image
docker rmi -f $SOFT_IMAGE_NAME:$SOFT_IMAGE_TAG

# build Dockerfile
docker build --build-arg GOPROXY=$GOPROXY -t $SOFT_IMAGE_NAME:$SOFT_IMAGE_TAG .

# remove dangling images
docker image prune -f

docker images | grep ${SOFT_NAME}

# run a container to print the soft version
docker run --rm --name running-${SOFT_NAME}-version $SOFT_IMAGE_NAME:$SOFT_IMAGE_TAG ${SOFT_NAME} -v

# push the image to the DockerHub
# docker push $SOFT_IMAGE_NAME:$SOFT_IMAGE_TAG
