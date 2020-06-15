#! /bin/bash

set -e

dir_name="$(dirname $0)"

source ${dir_name}/common.sh

which yarn >/dev/null 2>&1
YARN_EXISTS=$?
if [ $YARN_EXISTS -ne 0 ]; then
  echo "yarn is missing..."
  echo "perhaps try 'brew install yarn'"
  exit -1
fi

which packr2 >/dev/null 2>&1
PACKR2_EXISTS=$?
if [ $PACKR2_EXISTS -ne 0 ]; then
  echo "gobuffalo/packr/v2 is missing..."
  echo "check github.com/gobuffalo/packr/v2 for installation directions"
  exit -1
fi

ARG1=${1:-"build"}

function execute() {
  if [ "$1" == "clean" ]; then
    pushd ${REPO_ROOT}/pkg/spa
    packr2 clean
    popd
  else
    pushd ${REPO_ROOT}/ui/
    yarn build
    popd

    pushd ${REPO_ROOT}/pkg/spa
    packr2
    popd
  fi
}

execute $ARG1
