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

var (
	isNoTypes bool
	// KeysCommand is the command to search keys with partial match
	KeysCommand = &cobra.Command{
		Use:     "keys [query]",
		Short:   "Search keys in your parameter store",
		Example: fmt.Sprintf(queryExample, "keys", "keys", "keys", "keys"),
		Args:    cobra.RangeArgs(0, 1),
		PreRunE: initializeCredential,
		RunE: func(cmd *cobra.Command, args []string) error {
			return keys(args)
		},
	}
)

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
		begin := strings.Index(key, query)
		end := begin + len(query)
		if !isNoColor {
			key = key[0:begin] + color.RedString(key[begin:end]) + key[end:]
		}
		if isNoTypes {
			fmt.Fprintf(w, "%s\n", key)
		} else {
			fmt.Fprintf(w, "%s\tType: %s\n", key, v.Type)
		}
	}
	w.Flush()

	return nil
}

func init() {
	KeysCommand.PersistentFlags().BoolVar(&isNoTypes, "no-types", false, "Turn off parameter type shows")
}
