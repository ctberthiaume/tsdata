#!/bin/bash
# Build tsdata command-line tool for 64-bit MacOS and Linux

[[ -d tsdata.darwin-amd64 ]] && rm -rf tsdata.darwin-amd64
[[ -d tsdata.linux-amd64 ]] && rm -rf tsdata.linux-amd64
GOOS=darwin GOARCH=amd64 go build -o tsdata.darwin-amd64/tsdata cmd/tsdata/main.go || exit 1
GOOS=linux GOARCH=amd64 go build -o tsdata.linux-amd64/tsdata cmd/tsdata/main.go || exit 1
zip -q -r tsdata.darwin-amd64.zip tsdata.darwin-amd64 || exit 1
zip -q -r tsdata.linux-amd64.zip tsdata.linux-amd64 || exit 1
