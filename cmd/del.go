package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	isForceDelete bool
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

	param, err := ssmManager.GetParameter(key, false)
	if err != nil {
		fmt.Fprintln(ErrWriter, color.YellowString(fmt.Sprintf("WARN: %s is not found. Nothing to do.", key)))
		return nil
	}

	if !isForceDelete {
		fmt.Fprintf(ErrWriter, "Delete %s (value: %s) ? (Y/n)\n", key, param.Value)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			yn := scanner.Text()

			if yn == "Y" || yn == "y" {
				break
			} else if yn == "N" || yn == "n" {
				return nil
			} else {
				fmt.Println("(Y/n) ?")
			}
		}
	}

	if err := ssmManager.DeleteParameter(key); err != nil {
		fmt.Fprintln(ErrWriter, color.RedString(ErrMsgDeleteParameter))
		fmt.Fprintln(ErrWriter, err)
		return err
	}

	return nil
}

func init() {
	DelCommand.PersistentFlags().BoolVarP(&isForceDelete, "force", "f", false, "force to delete key and value")
}
