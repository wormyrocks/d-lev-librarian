package main

import (
	"fmt"
	"github.com/webview/webview"
)

func webui_init() {
	fmt.Print("hello world")
	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("D-Lev Librarian")
	w.SetSize(640, 480, webview.HintNone)
	w.SetHtml("<pre>hello</pre>")
	w.Run()
}
