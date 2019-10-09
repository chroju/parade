package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/chroju/para/cmd"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "para",
		Short: "SSM parameters",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Usage: keys")
		},
	}

	rootCmd.AddCommand(cmd.KeysCommand)
	rootCmd.AddCommand(cmd.GetCommand)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}