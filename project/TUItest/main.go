package main

import (
	"fmt"
	"os"

	"TUItest/tui"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "tui" {
		tui.RunTUI()
		return
	}

	fmt.Println("Use 'tui' command to launch the TUI interface")
}
