package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/chroju/para/ssmctl"
)

var GetCommand = &cobra.Command{
	Use: "get",
	Short: "Get key value",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ssmManager, err := ssmctl.New()
		if err != nil {
			fmt.Println(err)
		}

		resp, err := ssmManager.GetParameters(args[0])
		if err != nil {
			fmt.Println(err)
		}

		for _, v := range resp {
			fmt.Printf("%v\n", *v.Value)
		}
	},
}