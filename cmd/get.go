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

type getOption struct {
	Query        string
	Option       string
	IsDecryption bool
	IsNoColor    bool

	SSMManager ssmctl.SSMManager

	Out    io.Writer
	ErrOut io.Writer
}

func newGetCommand(globalOption *GlobalOption) *cobra.Command {
	o := &getOption{}

	cmd := &cobra.Command{
		Use:     "get <key>",
		Short:   "Get the value of specified key in your parameter store.",
		Example: queryExampleGet,
		Args:    cobra.ExactArgs(1),
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
			return o.get()
		},
	}

	cmd.PersistentFlags().BoolVarP(&o.IsDecryption, "decrypt", "d", false, "Get the value by decrypting it")
	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	return cmd
}

func (o *getOption) get() error {
	resp, err := o.SSMManager.DescribeParameters(o.Query, o.Option)
	if err != nil {
		return fmt.Errorf("%s\n%s", ErrMsgDescribeParameters, err)
	}

	if o.Option == ssmctl.DescribeOptionEquals {
		resp, err := o.SSMManager.GetParameter(o.Query, o.IsDecryption)
		if err != nil {
			if strings.Contains(err.Error(), "ParameterNotFound") {
				return nil
			}
			return err
		}
		fmt.Fprintln(o.Out, resp.Value)
		return nil
	}

	w := tabwriter.NewWriter(o.Out, 0, 2, 2, ' ', 0)

	for _, v := range resp {
		if err = o.getAndPrintParameter(w, v.Name); err != nil {
			return fmt.Errorf("%s\n%s", ErrMsgGetParameter, err)
		}
	}
	w.Flush()

	return nil
}

func (o *getOption) getAndPrintParameter(w *tabwriter.Writer, parameter string) error {
	resp, err := o.SSMManager.GetParameter(parameter, o.IsDecryption)
	if err != nil {
		return err
	}

	replacedLF := "\\n"
	begin := strings.Index(parameter, o.Query)
	end := begin + len(o.Query)
	if !o.IsNoColor {
		parameter = parameter[0:begin] + color.RedString(parameter[begin:end]) + parameter[end:]
		replacedLF = color.YellowString("\\n")
	}
	value := strings.ReplaceAll(resp.Value, "\n", replacedLF)
	fmt.Fprintf(w, "%s\t%s\n", parameter, value)

	return nil
}
