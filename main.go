package main

import (
	"os"
	"github.com/chroju/para/cmd"
)

func main() {
	cmd.Execute(os.Stdout, os.Stderr)
}