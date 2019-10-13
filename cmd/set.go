package cmd

import (
	"fmt"
	"os"

	"github.com/chroju/parade/ssmctl"
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
		Run: func(cmd *cobra.Command, args []string) {
			set(args)
		},
	}
)

func set(args []string) {
	key := args[0]
	value := args[1]

	ssmManager, err := ssmctl.New()
	if err != nil {
		fmt.Fprintln(ErrWriter, err)
		os.Exit(1)
	}

	if err = ssmManager.PutParameter(key, value, isEncryption, isForce); err != nil {
		fmt.Fprintln(ErrWriter, err)
		os.Exit(1)
	}

	fmt.Fprintln(ErrWriter, "done.")
}

func init() {
	SetCommand.PersistentFlags().BoolVarP(&isEncryption, "encrypt", "e", false, "set value with encryption")
	SetCommand.PersistentFlags().BoolVarP(&isForce, "force", "f", false, "force to overwrite the existing value")
}
