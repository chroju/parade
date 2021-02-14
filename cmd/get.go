package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	isAmbiguous  bool
	isDecryption bool

	// GetCommand is the command to get values of the specified keys
	GetCommand = &cobra.Command{
		Use:   "get",
		Short: "Get key value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return get(args)
		},
	}
)

func get(args []string) error {
	w := tabwriter.NewWriter(StdWriter, 0, 2, 2, ' ', 0)
	query := args[0]

	if isAmbiguous {
		resp, err := ssmManager.DescribeParameters(query)
		if err != nil {
			fmt.Fprintln(ErrWriter, color.RedString(ErrMsgDescribeParameters))
			fmt.Fprintln(ErrWriter, err)
			return err
		}

		for _, v := range resp {
			index := strings.Index(v.Name, query)
			if err = getAndPrintParameter(w, v.Name, index, index+len(query)); err != nil {
				fmt.Fprintln(ErrWriter, color.RedString(ErrMsgGetParameter))
				fmt.Fprintln(ErrWriter, err)
				break
			}
		}
	} else {
		getAndPrintParameter(w, query, 0, 0)
		if err := getAndPrintParameter(w, query, 0, 0); err != nil {
			fmt.Fprintln(ErrWriter, color.RedString(ErrMsgGetParameter))
			fmt.Fprintln(ErrWriter, err)
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

	coloredKey := key[0:begin] + color.RedString(key[begin:end]) + key[end:]
	value := strings.ReplaceAll(resp.Value, "\n", color.YellowString("\\n"))
	fmt.Fprintf(w, "%s\t%s\n", coloredKey, value)

	return nil
}

func init() {
	GetCommand.PersistentFlags().BoolVarP(&isAmbiguous, "ambiguous", "a", false, "Get all values that partially match the specified key")
	GetCommand.PersistentFlags().BoolVarP(&isDecryption, "decrypt", "d", false, "Get the value by decrypting it")
}
