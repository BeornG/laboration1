package main

import (
	"diagra/cmd/cli"
	"diagra/cmd/tui"

	"os"
)

func main() {
	input := len(os.Args)

	if input > 1 {
		cli.RunCLI(os.Args[1:])
	} else {
		tui.RunTUI()
	}

}
