package main

import (
	"os"

	"github.com/arunsanna/msgcli/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
