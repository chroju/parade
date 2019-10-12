package cmd

import (
	"fmt"
	"os"
	"strings"
	"github.com/spf13/cobra"
	"github.com/chroju/para/ssmctl"
	"github.com/fatih/color"
	"text/tabwriter"
)

var (
	isAmbiguous bool

	GetCommand = &cobra.Command{
		Use: "get",
		Short: "Get key value",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
			query := args[0]
			ssmManager, err := ssmctl.New()
			if err != nil {
				fmt.Println(err)
			}

			if isAmbiguous {
				resp, err := ssmManager.DescribeParameters()
				if err != nil {
					fmt.Println(err)
				}

				for _, v := range resp {
					index := strings.Index(*v.Name, query)
					if index >= 0 {
						resp, err := ssmManager.GetParameters(v)
						if err != nil {
							fmt.Println(err)
						}
						printValuesWithColor(w, ssmManager, *v.Name, *resp[0].Value, index, index+len(query))
					}
				}
			} else {
				resp, err := ssmManager.GetParameters(query)
				if err != nil {
					fmt.Println(err)
				}
				printValue(w, query, *resp[0].Value)
			}
		},
	}
)

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
	GetCommand.PersistentFlags().BoolVarP(&isAmbiguous, "ambiguous", "a", false, "get all keys")
}