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

which pkger >/dev/null 2>&1
PKGER_EXISTS=$?
if [ $PKGER_EXISTS -ne 0 ]; then
  echo "markbates/pkger is missing..."
  echo "check https://github.com/markbates/pkger for installation directions"
  exit -1
fi

pushd ${REPO_ROOT}/ui/
yarn build
popd

pushd ${REPO_ROOT}
go generate ./...
popd
