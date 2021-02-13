package cmd

import (
	"fmt"
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
	ssmManager, err := ssmctl.New()
	if err != nil {
		fmt.Fprintln(ErrWriter, err)
	}

	if isAmbiguous {
		resp, err := ssmManager.DescribeParameters()
		if err != nil {
			fmt.Fprintln(ErrWriter, err)
		}

		for _, v := range resp {
			index := strings.Index(v.Name, query)
			if index >= 0 {
				resp, err := ssmManager.GetParameter(v.Name, isDecryption)
				if err != nil {
					fmt.Fprintln(ErrWriter, err)
				}
				printValuesWithColor(w, ssmManager, v.Name, resp.Value, index, index+len(query))
			}
		}
	} else {
		resp, err := ssmManager.GetParameter(query, isDecryption)
		if err != nil {
			return
		}
		printValue(w, query, resp.Value)
	}
}

func printValue(w *tabwriter.Writer, key string, value string) {
	fmt.Fprintf(w, "%s\t%s\n", key, value)
	w.Flush()
}

func printValuesWithColor(w *tabwriter.Writer, s *ssmctl.SSMManager, key string, value string, begin int, end int) {
	coloredKey := key[0:begin] + color.RedString(key[begin:end]) + key[end:]
	printValue(w, coloredKey, value)
	w.Flush()
}

func init() {
	GetCommand.PersistentFlags().BoolVarP(&isAmbiguous, "ambiguous", "a", false, "get all values of the keys partial match")
	GetCommand.PersistentFlags().BoolVarP(&isDecryption, "decrypt", "d", false, "get keys with decription")
}
