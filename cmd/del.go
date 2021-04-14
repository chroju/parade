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

type delOption struct {
	Key           string
	IsForceDelete bool

	SSMManager ssmctl.SSMManager

	Out    io.Writer
	ErrOut io.Writer
}

func newDelCommand(globalOption *GlobalOption) *cobra.Command {
	o := &delOption{}

	cmd := &cobra.Command{
		Use:   "del <key>",
		Short: "Delete key and value in your parameter store.",
		Args:  cobra.ExactArgs(1),
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

			o.Out = globalOption.Out
			o.ErrOut = globalOption.ErrOut

			return o.del()
		},
	}

	cmd.PersistentFlags().BoolVarP(&o.IsForceDelete, "force", "f", false, "Force deletion of key and value\nDefault, display a prompt to confirm that you want to delete")
	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	return cmd
}

func (o *delOption) del() error {
	param, err := o.SSMManager.GetParameter(o.Key, false)
	if err != nil {
		fmt.Fprintln(o.ErrOut, color.YellowString(fmt.Sprintf("WARN: `%s` is not found. Nothing to do.", o.Key)))
		return nil
	}

	if !o.IsForceDelete {
		fmt.Fprintf(o.ErrOut, "Delete `%s` (value: %s) ? (Y/n)\n", o.Key, param.Value)
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

	if err := o.SSMManager.DeleteParameter(o.Key); err != nil {
		return fmt.Errorf("%s\n%s", ErrMsgDeleteParameter, err)
	}

	return nil
}
