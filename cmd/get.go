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
		Run: func(cmd *cobra.Command, args []string) {
			get(args)
		},
	}
)

func get(args []string) {
	w := tabwriter.NewWriter(StdWriter, 0, 2, 2, ' ', 0)
	query := args[0]

	if isAmbiguous {
		resp, err := ssmManager.DescribeParameters(query)
		if err != nil {
			fmt.Fprintln(ErrWriter, err)
		}

		for _, v := range resp {
			index := strings.Index(v.Name, query)
			getAndPrintParameter(w, v.Name, index, index+len(query))
		}
	} else {
		getAndPrintParameter(w, query, 0, 0)
	}
	w.Flush()
}

func getAndPrintParameter(w *tabwriter.Writer, key string, begin int, end int) {
	resp, err := ssmManager.GetParameter(key, isDecryption)
	if err != nil {
		fmt.Fprintln(ErrWriter, err)
	}
	coloredKey := key[0:begin] + color.RedString(key[begin:end]) + key[end:]
	fmt.Fprintf(w, "%s\t%s\n", coloredKey, resp.Value)
}

func init() {
	GetCommand.PersistentFlags().BoolVarP(&isAmbiguous, "ambiguous", "a", false, "get all values of the keys partial match")
	GetCommand.PersistentFlags().BoolVarP(&isDecryption, "decrypt", "d", false, "get keys with decription")
}
