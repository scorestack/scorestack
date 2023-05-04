#!/bin/bash
set -euxo pipefail

export PATH="$PATH:$GOPATH/bin"
cd $HOME/scorestack/dynamicbeat
make
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2
golangci-lint run -v
