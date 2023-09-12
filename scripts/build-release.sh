#!/usr/bin/env bash

echo "current branch is $(git branch --show-current)"

# update repository
#git pull --no-rebase

# update the last git commit
echo -e "$(git rev-parse HEAD)\c" >internal/version/commit

# set GOPROXY environment variable
# export GOPROXY=https://goproxy.cn

export SOFT_RELEASE_GO_VERSION
export SOFT_RELEASE_VERSION
export SOFT_NAME="gofs"
export SOFT_PREFIX="${SOFT_NAME}_"

function init_version {
  go build -v -o . ./...

  SOFT_RELEASE_GO_VERSION=$(go version | awk '{print $3}')
  SOFT_RELEASE_VERSION=$(./${SOFT_NAME} -v | awk 'NR==1 {print $3}')
}

function build_release {
  # release path, for example, gofs_go1.21.1_arm64_linux_v0.8.0
  SOFT_RELEASE="${SOFT_PREFIX}${SOFT_RELEASE_GO_VERSION}_${GOARCH}_${GOOS}_${SOFT_RELEASE_VERSION}"

  rm -rf "$SOFT_RELEASE"
  mkdir "$SOFT_RELEASE"

  # build
  go build -v -o . ./...

  if [ "$GOOS" == "windows" ]; then
    go build -v -ldflags="-H windowsgui" -o ./${SOFT_NAME}_background.exe ./cmd/${SOFT_NAME}
    mv ${SOFT_NAME}.exe ${SOFT_NAME}_background.exe "$SOFT_RELEASE/"
    # windows release archive
    zip -r "$SOFT_RELEASE.zip" "$SOFT_RELEASE"
  else
    mv ${SOFT_NAME} "$SOFT_RELEASE/"
    # release archive
    tar -zcvf "$SOFT_RELEASE.tar.gz" "$SOFT_RELEASE"
  fi
  rm -rf "$SOFT_RELEASE"
}

init_version

############################## linux-amd64-release ##############################

# set go env
export GOOS=linux
export GOARCH=amd64

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

build_release

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

ls -alh | grep ${SOFT_PREFIX}
