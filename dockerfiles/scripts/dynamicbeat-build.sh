#!/bin/bash
set -euxo pipefail

cd $GOPATH/src/github.com/s-newman/scorestack/dynamicbeat
ls -l $GOPATH
ls -l $GOPATH/src
ls -l $GOPATH/src/github.com
ls -l $GOPATH/src/github.com/s-newman
ls -l $GOPATH/src/github.com/s-newman/scorestack
ls -l $GOPATH/src/github.com/s-newman/scorestack/dynamicbeat
go get
make