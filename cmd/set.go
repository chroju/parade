package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	isEncryption bool
	isForce      bool

	// SetCommand is the command to set key value
	SetCommand = &cobra.Command{
		Use:     "set <key> <value>",
		Short:   "Set key and value in your parameter store. Overwriting is also possible.",
		Args:    cobra.ExactArgs(2),
		PreRunE: initializeCredential,
		RunE: func(cmd *cobra.Command, args []string) error {
			outWriter := os.Stdout
			errWriter := os.Stderr
			return set(args, outWriter, errWriter)
		},
	}
)

func set(args []string, outWriter, errWriter io.Writer) error {
	key := args[0]
	value := args[1]

	param, err := ssmManager.GetParameter(key, false)
	if err == nil && !isForce {
		fmt.Fprintf(errWriter, color.YellowString(fmt.Sprintf("WARN: `%s` already exists.\n", key)))
		fmt.Fprintf(errWriter, "Overwrite `%s` (value: %s) ? (Y/n)\n", key, param.Value)
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

	if err := ssmManager.PutParameter(key, value, isEncryption, isForce); err != nil {
		return fmt.Errorf("%s\n%s", ErrMsgPutParameter, err)
	}

	return nil
}

func init() {
	SetCommand.PersistentFlags().BoolVarP(&isEncryption, "encrypt", "e", false, "Encrypt the value and set it")
	SetCommand.PersistentFlags().BoolVarP(&isForce, "force", "f", false, "Force overwriting of existing values\nDefault, display a prompt to confirm that\nyou want to overwrite if the specified key already exists")
}
