package cmd

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/chroju/parade/ssmctl"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type keysOption struct {
	Query     string
	Option    string
	IsNoTypes bool
	IsNoColor bool

	SSMManager ssmctl.SSMManager

	Out    io.Writer
	ErrOut io.Writer
}

func newKeysCommand(globalOption *GlobalOption) *cobra.Command {
	o := &keysOption{}

	cmd := &cobra.Command{
		Use:     "keys [query]",
		Short:   "Search and show keys in your parameter store.",
		Example: queryExampleKeys,
		Args:    cobra.RangeArgs(0, 1),
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
			query := ""
			if len(args) != 0 {
				query = args[0]
			}
			query, option, err := queryParser(query)
			if err != nil {
				return err
			}
			o.Query = query
			o.Option = option
			o.IsNoColor = globalOption.IsNoColor

			o.Out = globalOption.Out
			o.ErrOut = globalOption.ErrOut

			return o.keys()
		},
	}

	cmd.Flags().BoolVar(&o.IsNoTypes, "no-types", false, "Turn off parameter type shows")

	return cmd
}

func (o *keysOption) keys() error {
	resp, err := o.SSMManager.DescribeParameters(o.Query, o.Option)
	if err != nil {
		return fmt.Errorf("%s\n%s", ErrMsgDescribeParameters, err)
	}

	w := tabwriter.NewWriter(o.Out, 0, 2, 2, ' ', 0)
	for _, v := range resp {
		key := v.Name
		begin := strings.Index(key, o.Query)
		end := begin + len(o.Query)
		if !o.IsNoColor {
			key = key[0:begin] + color.RedString(key[begin:end]) + key[end:]
		}
		if o.IsNoTypes {
			fmt.Fprintf(w, "%s\n", key)
		} else {
			fmt.Fprintf(w, "%s\tType: %s\n", key, v.Type)
		}
	}
	w.Flush()

	return nil
}
