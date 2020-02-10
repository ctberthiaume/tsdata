#!/bin/bash
# Build tsdata command-line tool for 64-bit MacOS and Linux

VERSION=$(git describe --long --dirty)
GOOS=darwin GOARCH=amd64 go build -o "tsdata.${VERSION}.darwin-amd64/tsdata" cmd/tsdata/main.go || exit 1
GOOS=linux GOARCH=amd64 go build -o "tsdata.${VERSION}.linux-amd64/tsdata" cmd/tsdata/main.go || exit 1
zip -q -r "tsdata.${VERSION}.darwin-amd64.zip" "tsdata.${VERSION}.darwin-amd64" || exit 1
zip -q -r "tsdata.${VERSION}.linux-amd64.zip" "tsdata.${VERSION}.linux-amd64"|| exit 1
