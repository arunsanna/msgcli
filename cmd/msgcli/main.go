package main

import (
	"os"

	"github.com/skylarbpayne/msgcli/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
