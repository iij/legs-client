#!/usr/bin/env bash

git pull --tags

VERSION=$(git tag -l --sort=v:refname | tail -n 1)
HASH=$(git rev-parse --verify HEAD)
BUILDDATE=$(date '+%Y/%m/%d %H:%M:%S %Z')
GOVERSION=$(go version)

echo version: $VERSION
echo hash: $HASH
echo build date: $BUILDDATE
echo go version: $GOVERSION

gox -ldflags="-s -w -X main.version=$VERSION -X main.hash=$HASH -X \"main.builddate=${BUILDDATE}\" -X \"main.goversion=$GOVERSION\"" \
	-os="linux netbsd darwin" \
	-output="./dist/legsc_{{.OS}}_{{.Arch}}" \
	./...
