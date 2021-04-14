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

type setOption struct {
	Key          string
	Value        string
	IsEncryption bool
	IsForce      bool

	SSMManager ssmctl.SSMManager

	Out    io.Writer
	ErrOut io.Writer
}

func newSetCommand(globalOption *GlobalOption) *cobra.Command {
	o := &setOption{}

	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set key and value in your parameter store. Overwriting is also possible.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.SSMManager = globalOption.SSMManager
			if o.SSMManager == nil {
				ssmManager, err := initializeCredential(globalOption.Profile, globalOption.Region)
				if err != nil {
					return err
				}
				o.SSMManager = ssmManager
			}

			args = cmd.Flags().Args()
			o.Key = args[0]
			o.Value = args[1]

			o.Out = globalOption.Out
			o.ErrOut = globalOption.ErrOut
			return o.set()
		},
	}

	cmd.PersistentFlags().BoolVarP(&o.IsEncryption, "encrypt", "e", false, "Encrypt the value and set it")
	cmd.PersistentFlags().BoolVarP(&o.IsForce, "force", "f", false, "Force overwriting of existing values\nDefault, display a prompt to confirm that\nyou want to overwrite if the specified key already exists")
	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	return cmd
}

func (o *setOption) set() error {
	param, err := o.SSMManager.GetParameter(o.Key, false)
	if err == nil && !o.IsForce {
		fmt.Fprintf(o.ErrOut, color.YellowString(fmt.Sprintf("WARN: `%s` already exists.\n", o.Key)))
		fmt.Fprintf(o.ErrOut, "Overwrite `%s` (value: %s) ? (Y/n)\n", o.Key, param.Value)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			yn := scanner.Text()

			if yn == "Y" || yn == "y" {
				break
			} else if yn == "N" || yn == "n" {
				return nil
			} else {
				fmt.Fprint(o.ErrOut, "(Y/n) ?")
			}
		}
	}

	if err := o.SSMManager.PutParameter(o.Key, o.Value, o.IsEncryption, o.IsForce); err != nil {
		return fmt.Errorf("%s\n%s", ErrMsgPutParameter, err)
	}

	return nil
}
