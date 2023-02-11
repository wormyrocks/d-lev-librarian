//go:build webui

package main

import (
	"embed"
	"fmt"
	"github.com/webview/webview"
)

const do_webview_debug = true


/*
  We don't want to run a local web server because of Windows permissions issues:
  https://github.com/webview/webview/issues/556#issuecomment-805672457

  And that makes it annoying to include a bunch of different files that all reference
  each other, because the backend then needs to serve multiple.

  There are probably ways around this, but the easiest is probably to package the
  whole shebang as a single HTML file.
*/

//go:embed webui/build/compiled_webapp.html
var http_payload embed.FS

func addTwo_(s int) int {
	fmt.Println(s, " + 2 = ", s + 2)
	return s + 2
}

func start_ui() {
	fmt.Println("hello world")
	w := webview.New(do_webview_debug)
	defer w.Destroy()
	w.SetTitle("D-Lev Librarian")
	w.Bind("addTwo", addTwo_)
	w.SetSize(640, 480, webview.HintFixed)
	data, _ := http_payload.ReadFile("webui/build/compiled_webapp.html")
	w.SetHtml(string(data))
	w.Run()
}
