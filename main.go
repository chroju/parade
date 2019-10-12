package main

import (
	"os"
	"github.com/chroju/parade/cmd"
)

func main() {
	if err := cmd.Execute(os.Stdout, os.Stderr); err != nil {
		os.Exit(1)
	}
}