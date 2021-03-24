#!/bin/bash
set -euxo pipefail

export PATH="$PATH:$GOPATH/bin"
cd $HOME/scorestack/dynamicbeat
make