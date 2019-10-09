package cmd

import (
	"fmt"
	"strings"
	"github.com/spf13/cobra"
	"github.com/chroju/para/ssmctl"
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
			if strings.Contains(*v.Name, args[0]) {
				fmt.Printf("%v\n", *v.Name)
			}
		}
	},
}