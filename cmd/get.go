package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/chroju/parade/ssmctl"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	isAmbiguous  bool
	isDecryption bool

	// GetCommand is the command to get values of the specified keys
	GetCommand = &cobra.Command{
		Use:     "get <key>",
		Short:   "Get the value of specified key in your parameter store.",
		Example: queryExampleGet,
		Args:    cobra.ExactArgs(1),
		PreRunE: initializeCredential,
		RunE: func(cmd *cobra.Command, args []string) error {
			outWriter := os.Stdout
			errWriter := os.Stderr
			return get(args, outWriter, errWriter)
		},
	}
)

func get(args []string, outWriter, errWiter io.Writer) error {
	w := tabwriter.NewWriter(outWriter, 0, 2, 2, ' ', 0)
	query, option, err := queryParser(args[0])
	if err != nil {
		return err
	}

	if option == ssmctl.DescribeOptionEquals {
		resp, err := ssmManager.GetParameter(query, isDecryption)
		if err != nil {
			return err
		}
		fmt.Fprintln(outWriter, resp.Value)
		return nil
	}

	resp, err := ssmManager.DescribeParameters(query, option)
	if err != nil {
		return fmt.Errorf("%s\n%s", ErrMsgDescribeParameters, err)
	}

	for _, v := range resp {
		index := strings.Index(v.Name, query)
		if err = getAndPrintParameter(w, v.Name, index, index+len(query)); err != nil {
			return fmt.Errorf("%s\n%s", ErrMsgGetParameter, err)
		}
	}
	w.Flush()

	return nil
}

func getAndPrintParameter(w *tabwriter.Writer, key string, begin int, end int) error {
	resp, err := ssmManager.GetParameter(key, isDecryption)
	if err != nil {
		return err
	}

	replacedLF := "\\n"
	if !isNoColor {
		key = key[0:begin] + color.RedString(key[begin:end]) + key[end:]
		replacedLF = color.YellowString("\\n")
	}
	value := strings.ReplaceAll(resp.Value, "\n", replacedLF)
	fmt.Fprintf(w, "%s\t%s\n", key, value)

	return nil
}

func init() {
	GetCommand.PersistentFlags().BoolVarP(&isDecryption, "decrypt", "d", false, "Get the value by decrypting it")
}
