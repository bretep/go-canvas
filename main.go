package main

import (
	"fmt"
	"os"

	"github.com/harrybrwn/go-canvas/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
