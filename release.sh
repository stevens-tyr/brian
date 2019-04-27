#!/bin/bash
GIT_HASH=$(git rev-parse --short HEAD)
echo Compiling Release $GIT_HASH
rm -rf brian
GOOS=linux GOARCH=amd64 go build -o brian
echo Creating Release $GIT_HASH
hub release create -m "version: $GIT_HASH" -a brian $GIT_HASH
