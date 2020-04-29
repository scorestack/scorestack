#!/bin/bash
set -euxo pipefail

cd $GOPATH/src/github.com/s-newman/scorestack/dynamicbeat
make testsuite