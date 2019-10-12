package cmd

import (
	"fmt"
	"strings"
	"github.com/spf13/cobra"
	"github.com/chroju/para/ssmctl"
	"github.com/fatih/color"
)

var KeysCommand = &cobra.Command{
	Use: "keys",
	Short: "Get keys",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ssmManager, err := ssmctl.New()
		if err != nil {
			fmt.Println(err)
		}

		resp, err := ssmManager.DescribeParameters()
		if err != nil {
			fmt.Println(err)
		}

		for _, v := range resp {
			key := *v.Name
			index := strings.Index(key, args[0])
			if index >= 0 {
				end := index + len(args[0])
				coloredKey := key[0:index] + color.RedString(key[index:end]) + key[end:]
				fmt.Printf("%v\n", coloredKey)
			}
		}
	},
}