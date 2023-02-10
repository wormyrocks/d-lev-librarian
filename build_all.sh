#! /bin/bash
GOOS=linux GOARCH=amd64 go build -o d-lin
GOOS=windows GOARCH=amd64 go build -o d-win.exe
GOOS=darwin GOARCH=amd64 go build -o d-mac
GOOS=darwin GOARCH=arm64 go build -o d-mm1
exit
