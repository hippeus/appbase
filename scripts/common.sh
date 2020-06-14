#! /bin/bash

REPO_ROOT=$(git rev-parse --show-toplevel)
GO_MOD_ROOT=$(go list -m)
APP_NAME=${1:-demoapp}
GOOS=${GOOS:=darwin}
GOARCH=${GOARCH:=amd64}
