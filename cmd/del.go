package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// DelCommand is the command to delete key value
	DelCommand = &cobra.Command{
		Use:   "del",
		Short: "Delete key value",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			del(args)
		},
	}
)

func del(args []string) {
	key := args[0]

	if err := ssmManager.DeleteParameter(key); err != nil {
		fmt.Fprintln(ErrWriter, err)
		os.Exit(1)
	}

	fmt.Fprintln(ErrWriter, "done.")
}

func init() {
}
