#!/usr/bin/env bash

echo "current branch is $(git branch --show-current)"

# update repository
#git pull --no-rebase

# update the last git commit
echo -e "$(git rev-parse HEAD)\c" >internal/version/commit

# set GOPROXY environment variable
# export GOPROXY=https://goproxy.cn

function build_release {
  # build
  go build -v -o . ./...

  # release path, for example, gofs_go1.20.1_arm64_linux_v0.6.0
  GOFS_RELEASE="gofs_${GOFS_RELEASE_GO_VERSION}_${GOARCH}_${GOOS}_${GOFS_RELEASE_VERSION}"

  rm -rf "$GOFS_RELEASE"
  mkdir "$GOFS_RELEASE"
  mv gofs "$GOFS_RELEASE/"

  # release archive
  tar -zcvf "$GOFS_RELEASE.tar.gz" "$GOFS_RELEASE"

  rm -rf "$GOFS_RELEASE"
}

############################## linux-amd64-release ##############################

# set go env
export GOOS=linux
export GOARCH=amd64

# build
go build -v -o . ./...

GOFS_RELEASE_GO_VERSION=$(go version | awk '{print $3}')
GOFS_RELEASE_VERSION=$(./gofs -v | awk 'NR==1 {print $3}')

build_release

############################## linux-amd64-release ##############################

############################## linux-arm64-release ##############################

export GOOS=linux
export GOARCH=arm64

build_release

############################## linux-arm64-release ##############################

############################# windows-release #############################

# set go env
export GOOS=windows
export GOARCH=amd64

# build
go build -v -o . ./...

# build with -ldflags="-H windowsgui" flag
go build -v -ldflags="-H windowsgui" -o ./gofs_background.exe ./cmd/gofs

# release path, for example, gofs_go1.20.1_amd64_windows_v0.6.0
GOFS_RELEASE="gofs_${GOFS_RELEASE_GO_VERSION}_${GOARCH}_${GOOS}_${GOFS_RELEASE_VERSION}"

mkdir "$GOFS_RELEASE"
mv gofs.exe gofs_background.exe "$GOFS_RELEASE/"

# windows release archive
zip -r "$GOFS_RELEASE.zip" "$GOFS_RELEASE"

rm -rf "$GOFS_RELEASE"

############################# windows-release #############################

############################## macOS-amd64-release ##############################

export GOOS=darwin
export GOARCH=amd64

build_release

############################## macOS-amd64-release ##############################

############################## macOS-arm64-release ##############################

export GOOS=darwin
export GOARCH=arm64

build_release

############################## macOS-arm64-release ##############################

# reset commit file
echo -e "\c" >internal/version/commit

ls -alh | grep gofs_
