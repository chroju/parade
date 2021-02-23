package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/chroju/parade/ssmctl"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const ()

// KeysCommand is the command to search keys with partial match
var KeysCommand = &cobra.Command{
	Use:     "keys [query]",
	Short:   "Search keys in your parameter store",
	Example: fmt.Sprintf(queryExample, "keys", "keys", "keys", "keys"),
	Args:    cobra.RangeArgs(0, 1),
	PreRunE: initializeCredential,
	RunE: func(cmd *cobra.Command, args []string) error {
		return keys(args)
	},
}

func keys(args []string) error {
	query := ""
	option := ssmctl.DescribeOptionEquals
	if len(args) != 0 {
		var err error
		query, option, err = queryParser(args[0])
		if err != nil {
			return err
		}
	}

	resp, err := ssmManager.DescribeParameters(query, option)
	if err != nil {
		return fmt.Errorf("%s\n%s", ErrMsgDescribeParameters, err)
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

	return nil
}
