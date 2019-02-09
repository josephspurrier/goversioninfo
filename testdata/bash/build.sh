#!/bin/bash

# Exit on error.
set -e

# Copy the icon to test file closing.
cp $GOPATH/src/github.com/josephspurrier/goversioninfo/testdata/resource/icon.ico $GOPATH/src/github.com/josephspurrier/goversioninfo/testdata/resource/icon2.ico

# Test Windows 32.
cd $GOPATH/src/github.com/josephspurrier/goversioninfo/testdata/example32
GOOS=windows GOARCH=386 go generate
GOOS=windows GOARCH=386 go build
rm $GOPATH/src/github.com/josephspurrier/goversioninfo/testdata/resource/icon.ico

# Restore the icon.
cp $GOPATH/src/github.com/josephspurrier/goversioninfo/testdata/resource/icon2.ico $GOPATH/src/github.com/josephspurrier/goversioninfo/testdata/resource/icon.ico

# Test Windows 64.
cd $GOPATH/src/github.com/josephspurrier/goversioninfo/testdata/example64
GOOS=windows GOARCH=amd64 go generate
GOOS=windows GOARCH=amd64 go build
rm $GOPATH/src/github.com/josephspurrier/goversioninfo/testdata/resource/icon.ico

# Reset the icons.
mv $GOPATH/src/github.com/josephspurrier/goversioninfo/testdata/resource/icon2.ico $GOPATH/src/github.com/josephspurrier/goversioninfo/testdata/resource/icon.ico
cd $GOPATH/src/github.com/josephspurrier/goversioninfo