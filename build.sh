#!/bin/bash

export CGO_ENABLED=0

go get -a -ldflags '-s' -v

go build -ldflags '-s' -x
