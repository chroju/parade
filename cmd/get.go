package cmd

import (
	"fmt"
	"strings"
	"github.com/spf13/cobra"
	"github.com/chroju/para/ssmctl"
)

var isAmbiguous bool

var GetCommand = &cobra.Command{
	Use: "get",
	Short: "Get key value",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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
				if strings.Contains(*v.Name, args[0]) {
					printValue(ssmManager, *v.Name)
				}
			}
		} else {
			printValue(ssmManager, args[0])
		}
	},
}

func printValue(s *ssmctl.SSMManager, key string) {
	resp, err := s.GetParameters(key)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range resp {
		fmt.Printf("%v\n", *v.Value)
	}
}

func init() {
	GetCommand.PersistentFlags().BoolVarP(&isAmbiguous, "ambiguous", "a", false, "get all keys")
}