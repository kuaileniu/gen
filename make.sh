#!/bin/sh
# chmod 775 make.sh
go fmt ./...
rm -rf bin

go build -gcflags=-m  -ldflags "-s -w" -i -o bin/gen