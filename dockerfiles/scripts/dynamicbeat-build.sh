#!/bin/bash
set -euxo pipefail

echo $GOPATH
cd $GOPATH/src/github.com/s-newman/scorestack/dynamicbeat
go get
make