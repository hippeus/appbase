#! /bin/bash

set -e

dir_name="$(dirname $0)"

source ${dir_name}/common.sh

go build -ldflags \
  "-X '${GO_MOD_ROOT}/pkg/buildinfo.version=$(git describe --tags --always --dirty)' \
  -X  '${GO_MOD_ROOT}/pkg/buildinfo.commit=$(git rev-list -1 HEAD)' \
  -X  '${GO_MOD_ROOT}/pkg/buildinfo.buildTime=$(date -u '+%Y-%m-%d %H:%M:%S %Z')'" \
  -o ${REPO_ROOT}/bin/${GOOS}_${GOARCH}/${APP_NAME} ${REPO_ROOT}/cmd/${APP_NAME}

cp ${REPO_ROOT}/config/${APP_NAME}.yaml ${REPO_ROOT}/bin/${GOOS}_${GOARCH}/
