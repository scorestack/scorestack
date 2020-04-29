#!/bin/bash
set -euxo pipefail

cd $GOPATH/src/github.com/s-newman/scorestack/dynamicbeat
go get
go get github.com/kardianos/govendor
make testsuite