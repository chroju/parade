package main

import (
	"os"
	"github.com/chroju/para/cmd"
)

func main() {
	if err := cmd.Execute(os.Stdout, os.Stderr); err != nil {
		os.Exit(1)
	}
}