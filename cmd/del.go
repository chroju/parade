package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/chroju/parade/ssmctl"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	isForceDelete bool
	// DelCommand is the command to delete key value
	DelCommand = &cobra.Command{
		Use:   "del <key>",
		Short: "Delete key and value in your parameter store.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			outWriter := os.Stdout
			errWriter := os.Stderr
			ssmManager, err := initializeCredential(flagProfile, flagRegion)
			if err != nil {
				return err
			}
			return del(args, ssmManager, outWriter, errWriter)
		},
	}
)

func del(args []string, ssmManager ssmctl.SSMManager, outWriter, errWriter io.Writer) error {
	key := args[0]

	param, err := ssmManager.GetParameter(key, false)
	if err != nil {
		fmt.Fprintln(errWriter, color.YellowString(fmt.Sprintf("WARN: `%s` is not found. Nothing to do.", key)))
		return nil
	}

	if !isForceDelete {
		fmt.Fprintf(errWriter, "Delete `%s` (value: %s) ? (Y/n)\n", key, param.Value)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			yn := scanner.Text()

			if yn == "Y" || yn == "y" {
				break
			} else if yn == "N" || yn == "n" {
				return nil
			} else {
				fmt.Fprint(errWriter, "(Y/n) ?")
			}
		}
	}

	if err := ssmManager.DeleteParameter(key); err != nil {
		return fmt.Errorf("%s\n%s", ErrMsgDeleteParameter, err)
	}

	return nil
}

func init() {
	DelCommand.PersistentFlags().BoolVarP(&isForceDelete, "force", "f", false, "Force deletion of key and value\nDefault, display a prompt to confirm that you want to delete")
}
