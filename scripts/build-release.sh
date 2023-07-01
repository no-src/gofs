#!/usr/bin/env bash

# switch to branch main
git checkout main

# update repository
#git pull --no-rebase

# update the last git commit
echo -e "$(git rev-parse main)\c" >internal/version/commit

# set GOPROXY environment variable
# export GOPROXY=https://goproxy.cn

############################## linux-release ##############################

# set go env for linux
export GOOS=linux
export GOARCH=amd64

# build gofs
go build -v -o . ./...

GOFS_RELEASE_GO_VERSION=$(go version | awk '{print $3}')
GOFS_RELEASE_VERSION=$(./gofs -v | awk 'NR==1 {print $3}')

# release path, for example, gofs_go1.20.1_amd64_linux_v0.6.0
GOFS_RELEASE="gofs_${GOFS_RELEASE_GO_VERSION}_${GOARCH}_${GOOS}_${GOFS_RELEASE_VERSION}"

rm -rf "$GOFS_RELEASE"
mkdir "$GOFS_RELEASE"
mv gofs "$GOFS_RELEASE/"

# linux release archive
tar -zcvf "$GOFS_RELEASE.tar.gz" "$GOFS_RELEASE"

rm -rf "$GOFS_RELEASE"

############################## linux-release ##############################

############################# windows-release #############################

# set go env for windows
export GOOS=windows
export GOARCH=amd64

# build gofs
go build -v -o . ./...

# build gofs with -ldflags="-H windowsgui" flag
go build -v -ldflags="-H windowsgui" -o ./gofs_background.exe ./cmd/gofs

# release path, for example, gofs_go1.20.1_amd64_windows_v0.6.0
GOFS_RELEASE="gofs_${GOFS_RELEASE_GO_VERSION}_${GOARCH}_${GOOS}_${GOFS_RELEASE_VERSION}"

mkdir "$GOFS_RELEASE"
mv gofs.exe gofs_background.exe "$GOFS_RELEASE/"

# windows release archive
zip -r "$GOFS_RELEASE.zip" "$GOFS_RELEASE"

rm -rf "$GOFS_RELEASE"

############################# windows-release #############################

############################## macOS-release ##############################

# set go env for macOS
export GOOS=darwin
export GOARCH=amd64

# build gofs
go build -v -o . ./...

# release path, for example, gofs_go1.20.1_amd64_darwin_v0.6.0
GOFS_RELEASE="gofs_${GOFS_RELEASE_GO_VERSION}_${GOARCH}_${GOOS}_${GOFS_RELEASE_VERSION}"

rm -rf "$GOFS_RELEASE"
mkdir "$GOFS_RELEASE"
mv gofs "$GOFS_RELEASE/"

# macOS release archive
tar -zcvf "$GOFS_RELEASE.tar.gz" "$GOFS_RELEASE"

rm -rf "$GOFS_RELEASE"

############################## macOS-release ##############################

# reset commit file
echo -e "\c" >internal/version/commit
