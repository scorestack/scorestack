#!/bin/bash
set -euxo pipefail

export PATH="$PATH:$GOPATH/bin"
cd $GOPATH/src/github.com/scorestack/scorestack/dynamicbeat
go get
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.25.1
golangci-lint run -v