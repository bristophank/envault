package main

import (
	"fmt"
	"os"

	"github.com/envault/envault/cmd/envault/commands"
)

func main() {
	if err := commands.Root().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
