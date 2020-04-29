#!/bin/bash
cd $GOPATH/src/github.com/s-newman/scorestack/dynamicbeat
make setup
go get
make build