#!/bin/bash
set -euxo pipefail

export PATH="$PATH:$GOPATH/bin"
cd $GOPATH/src/github.com/scorestack/scorestack/dynamicbeat
go get
#make update
make
#make testsuite