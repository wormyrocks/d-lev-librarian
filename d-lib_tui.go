//go:build tui

package main

import (
	"fmt"
	"github.com/rivo/tview"
)

// https://github.com/rivo/tview/wiki/DropDown
func start_ui() {
	fmt.Println("Starting tui")
	app := tview.NewApplication()
	dropdown := tview.NewDropDown().
	SetLabel("Select an option (hit Enter): ").
	SetOptions([]string{"First", "Second", "Third", "Fourth", "Fifth"}, nil)
	if err := app.SetRoot(dropdown, true).SetFocus(dropdown).Run(); err != nil {
		panic(err)
	}
}

