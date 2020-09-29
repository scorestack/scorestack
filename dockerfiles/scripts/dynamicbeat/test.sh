#!/bin/bash
set -euxo pipefail

export PATH="$PATH:$GOPATH/bin"
cd $GOPATH/src/github.com/scorestack/scorestack/dynamicbeat
go get
go get github.com/kardianos/govendor
make setup
make
#make testsuite