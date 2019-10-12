package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "parade",
	Short: "simple SSM parameters CLI",
	Long: `Parade is a simple CLI tool for AWS SSM parameter store.
Easy to read and writer key values in your parameter store.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: keys")
	},
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(KeysCommand)
	rootCmd.AddCommand(GetCommand)
}