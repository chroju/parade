package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/chroju/parade/ssmctl"
	"github.com/spf13/cobra"
)

// VERSION is cli tool version
const VERSION = "0.2.0"

const (
	// ErrMsgAWSProfileNotValid is an error message to notify aws profile is not valid
	ErrMsgAWSProfileNotValid = `AWS credential is not valid.

Use the --profile and --region options, or set the access keys and region in the environment variables.
`
	// ErrMsgQueryFormat is an error message about query format
	ErrMsgQueryFormat = "The query format is not valid."
	// ErrMsgDescribeParameters is an error message about DescribeParameters API
	ErrMsgDescribeParameters = "Failed to execute DescribeParameters API."
	// ErrMsgGetParameter is an error message about GetParameter API
	ErrMsgGetParameter = "Failed to execute GetParameter API."
	// ErrMsgPutParameter is an error message about PutParameter API
	ErrMsgPutParameter = "Failed to execute PutParameter API."
	// ErrMsgDeleteParameter is an error message about DeleteParameter API
	ErrMsgDeleteParameter = "Failed to execute DeleteParameter API."

	queryExample = `  %s command supports exact match, forward match, and partial match.
  It usually searches for exact matches.

  $ parade %s /MyService/Test

  Use * as a postfix, the search will be done as a forward match.

  $ parade %s /MyService*

  Furthermore, also use * as a prefix, it becomes a partial match.

  $ parade %s *Test*

  If you do not specify any queries, display all keys.
`
)

var (
	// StdWriter is the io.Writer for standard output
	StdWriter io.Writer
	// ErrWriter is the io.Writer for error output
	ErrWriter io.Writer

	profile    string
	region     string
	ssmManager *ssmctl.SSMManager

	rootCmd = &cobra.Command{
		Use:     "parade",
		Short:   "simple SSM parameters CLI",
		Version: VERSION,
		Long: `Parade is a simple CLI tool for AWS SSM parameter store.
	Easy to read and writer key values in your parameter store.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Use subcommand: keys, get, set, del")
		},
	}
)

// Execute executes the root command
func Execute(w io.Writer, e io.Writer) error {
	StdWriter = w
	ErrWriter = e

	// aws-sdk-go does not support the AWS_DEFAULT_REGION environment variable
	region = os.Getenv("AWS_DEFAULT_REGION")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "AWS profile")
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "AWS region")

	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(KeysCommand, GetCommand, SetCommand, DelCommand)
}

func initializeCredential(cmd *cobra.Command, args []string) error {
	var err error
	ssmManager, err = ssmctl.New(profile, region)
	if err != nil {
		return fmt.Errorf(ErrMsgAWSProfileNotValid)
	}
	return nil
}
