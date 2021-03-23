#!/bin/bash
set -euxo pipefail

export PATH="$PATH:$GOPATH/bin"
cd $HOME/scorestack
make
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.25.1
golangci-lint run -v