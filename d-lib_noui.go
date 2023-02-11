//go:build (!webui && !tui && !imgui)

package main

import ("fmt")

func start_ui() {
	fmt.Println("d-lib was not compiled with a UI.")
}
