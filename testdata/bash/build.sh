#!/bin/bash

# Exit on error.
set -e

# Test Windows 32.
cd testdata/example32
GOOS=windows GOARCH=386 go generate
GOOS=windows GOARCH=386 go build
rm example32.exe
rm resource.syso
cd ../../

# Test Windows 64.
cd testdata/example64
GOOS=windows GOARCH=amd64 go generate
GOOS=windows GOARCH=amd64 go build
rm example64.exe
rm resource.syso
cd ../../

# Test Windows 64 with Go output file.
cd testdata/example64-gofile
GOOS=windows GOARCH=amd64 go generate
GOOS=windows GOARCH=amd64 go build
rm example64-gofile.exe
rm resource.syso
rm versioninfo.go
cd ../../