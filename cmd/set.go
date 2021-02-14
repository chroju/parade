package cmd

import (
	"bufio"
	"fmt"
	"os"

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

	param, err := ssmManager.GetParameter(key, false)
	if err == nil && !isForce {
		fmt.Fprintf(ErrWriter, color.YellowString(fmt.Sprintf("WARN: `%s` already exists.\n", key)))
		fmt.Fprintf(ErrWriter, "Overwrite `%s` (value: %s) ? (Y/n)\n", key, param.Value)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			yn := scanner.Text()

			if yn == "Y" || yn == "y" {
				break
			} else if yn == "N" || yn == "n" {
				return nil
			} else {
				fmt.Fprint(ErrWriter, "(Y/n) ?")
			}
		}
	}

	if err := ssmManager.PutParameter(key, value, isEncryption, isForce); err != nil {
		fmt.Fprintln(ErrWriter, color.RedString(ErrMsgPutParameter))
		fmt.Fprintln(ErrWriter, err)
		return err
	}

	return nil
}

func init() {
	SetCommand.PersistentFlags().BoolVarP(&isEncryption, "encrypt", "e", false, "Encrypt the value and set it")
	SetCommand.PersistentFlags().BoolVarP(&isForce, "force", "f", false, "Force overwriting of existing values")
}
