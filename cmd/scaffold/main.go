package main

import (
	"os"

	"github.com/apexsudo/common/cmd/scaffold/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
