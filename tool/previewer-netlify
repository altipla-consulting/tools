#!/bin/bash

set -eu


APP=$(basename $0)


mkdir -p tmp/bin
go build -o tmp/bin/$APP ./cmd/$APP
tmp/bin/$APP "$@"
