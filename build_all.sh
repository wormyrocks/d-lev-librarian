#! /bin/bash
mkdir -p bin
GOOS=linux GOARCH=amd64 go build -o bin/d-lin
GOOS=windows GOARCH=amd64 go build -o bin/d-win.exe
GOOS=darwin GOARCH=amd64 go build -o bin/d-mac
GOOS=darwin GOARCH=arm64 go build -o bin/d-mm1
exit
