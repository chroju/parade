package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// KeysCommand is the command to search keys with partial match
var KeysCommand = &cobra.Command{
	Use:   "keys",
	Short: "Get keys",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		keys(args)
	},
}

func keys(args []string) {
	query := ""
	if len(args) != 0 {
		query = args[0]
	}

	resp, err := ssmManager.DescribeParameters(query)
	if err != nil {
		fmt.Fprintln(ErrWriter, err)
	}

	w := tabwriter.NewWriter(StdWriter, 0, 2, 2, ' ', 0)
	for _, v := range resp {
		key := v.Name
		index := strings.Index(key, query)
		end := index + len(query)
		coloredKey := key[0:index] + color.RedString(key[index:end]) + key[end:]
		fmt.Fprintf(w, "%s\tType: %s\n", coloredKey, v.Type)
	}
	w.Flush()
}
