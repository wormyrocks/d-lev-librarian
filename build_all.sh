#! /bin/bash
GOOS=linux GOARCH=amd64 go build -o d-lin
GOOS=windows GOARCH=amd64 go build -o d-win.exe
GOOS=darwin GOARCH=amd64 go build -o d-mac
GOOS=darwin GOARCH=arm64 go build -o d-mm1
GOOS=linux GOARCH=arm64 go build -o d-arm
GOOS=linux GOARCH=arm go build -o d-a32
exit
