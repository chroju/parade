package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// DelCommand is the command to delete key value
	DelCommand = &cobra.Command{
		Use:   "del",
		Short: "Delete key value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return del(args)
		},
	}
)

func del(args []string) error {
	key := args[0]

	if err := ssmManager.DeleteParameter(key); err != nil {
		fmt.Fprintln(ErrWriter, color.RedString(ErrMsgDeleteParameter))
		fmt.Fprintln(ErrWriter, err)
		return err
	}

	return nil
}

func init() {
}
