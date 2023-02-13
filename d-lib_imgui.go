//go:build imgui

// This apparently requires tdm-gcc to compile on Windows:
// https://jmeubank.github.io/tdm-gcc/articles/2021-05/10.3.0-release

// https://github.com/AllenDang/giu/blob/master/README.md#build-windows-version-on-macoslinux
// go build -ldflags "-s -w -H=windowsgui -extldflags=-static"

// This takes a very long time to build for the first time, have to be patient

package main

import (
	"fmt"
	g "github.com/AllenDang/giu"
)

var content string

func im_loop() {
	g.SingleWindow().Layout(
		g.Label("Hello world from giu"),
		g.InputTextMultiline(&content).Size(g.Auto, g.Auto),
	)
}

// https://github.com/AllenDang/giu/blob/master/examples/helloworld/helloworld.go
func start_ui() {
	fmt.Println("Starting imgui")
	wnd := g.NewMasterWindow("Hello world", 400, 200, 0)
	wnd.Run(im_loop)
}

