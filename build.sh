#!/bin/sh
WHAT=$1
rm ./build/*
go build -o build/floody-minesweeper ./src
if [ "$WHAT" = "all" ]; then
    env GOOS=windows GOARCH=amd64 go build -o build/floody-minesweeper.exe ./src
    env GOOS=darwin GOARCH=amd64 go build -o build/floody-minesweeper_osx ./src
fi

if [ "$WHAT" = "run" ]; then
    ./build/floody-minesweeper
fi