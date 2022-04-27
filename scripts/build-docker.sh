#!/usr/bin/env bash

# update git repository
# git pull --no-rebase

# update the last git commit
echo -e "$(git rev-parse main)\c" >version/commit

# set GOPROXY environment variable
# export GOPROXY=https://goproxy.cn

# set the golang version by GOIMAGETAG environment variable
export GOIMAGETAG=latest

# set the gofs docker image name by GOFSIMAGENAME environment variable
export GOFSIMAGENAME=nosrc/gofs

# set the gofs docker image tag by GOFSIMAGETAG environment variable
export GOFSIMAGETAG=latest

# build gofs
docker run --rm -v "$PWD":/usr/src/gofs -w /usr/src/gofs -e GOPROXY=$GOPROXY golang:$GOIMAGETAG go build -v -tags netgo -o . ./...

# reset commit file
echo -e "\c" >version/commit

# remove the existing old image
docker rmi -f $GOFSIMAGENAME

# build Dockerfile
docker build -t $GOFSIMAGENAME:$GOFSIMAGETAG .

# run a container to print the gofs version
docker run -it --rm --name running-gofs-version $GOFSIMAGENAME:$GOFSIMAGETAG gofs -v

# clear gofs binary file
rm -rf gofs
