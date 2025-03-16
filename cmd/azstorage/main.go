package main

import (
	"fmt"
	"os"
)

func main() {

	var rootCommand = newRootCommand()
	var err = rootCommand.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
