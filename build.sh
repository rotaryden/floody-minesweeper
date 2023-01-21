#!/bin/sh
WHAT=$1
rm ./build/*
go build -o build/proxx ./src
if [ "$WHAT" = "all" ]; then
    env GOOS=windows GOARCH=amd64 go build -o build/proxx.exe ./src
    env GOOS=darwin GOARCH=amd64 go build -o build/proxx_osx ./src
fi

if [ "$WHAT" = "run" ]; then
    ./build/proxx
fi