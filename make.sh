#!/bin/sh
# chmod 775 make.sh

rm -rf bin

go build -ldflags "-s -w" -i -o bin/gen