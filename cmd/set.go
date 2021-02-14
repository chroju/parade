package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	isEncryption bool
	isForce      bool

	// SetCommand is the command to set key value
	SetCommand = &cobra.Command{
		Use:   "set",
		Short: "Set key value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return set(args)
		},
	}
)

func set(args []string) error {
	key := args[0]
	value := args[1]

	if err := ssmManager.PutParameter(key, value, isEncryption, isForce); err != nil {
		fmt.Fprintln(ErrWriter, color.RedString(ErrMsgPutParameter))
		fmt.Fprintln(ErrWriter, err)
		return err
	}

	return nil
}

func init() {
	SetCommand.PersistentFlags().BoolVarP(&isEncryption, "encrypt", "e", false, "set value with encryption")
	SetCommand.PersistentFlags().BoolVarP(&isForce, "force", "f", false, "force to overwrite the existing value")
}
