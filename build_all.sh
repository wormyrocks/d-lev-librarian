#! /bin/bash

cd $(dirname ${0})

go_build() {
    # Clean build
    rm -rf bin
    mkdir -p bin

    # Compile Go stuff
    GOOS=linux GOARCH=amd64 go build -o bin/d-lin
    GOOS=windows GOARCH=amd64 go build -o bin/d-win.exe
    GOOS=darwin GOARCH=amd64 go build -o bin/d-mac
    GOOS=darwin GOARCH=arm64 go build -o bin/d-mm1
}

# Static site build must happen before Go build, as the Go compiler will embed the entire static webpage
# into the binary.
webui/build_webui.sh && go_build
exit