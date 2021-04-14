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
const VERSION = "0.3.0"

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

	longDescription = `Parade is a simple CLI tool for AWS SSM parameter store.
Easy to read and writer key values in your parameter store.

Note:
  Parade requires your AWS IAM user authentication.
  The same authentication method as AWS CLI is available.
  Tools like aws-vault can be used as well.
`

	queryExampleKeys = `  keys command supports exact match, forward match, and partial match.
  It usually searches for exact matches.

  $ parade keys /MyService/Test

  Use * as a postfix, the search will be done as a forward match.

  $ parade keys /MyService*

  Furthermore, also use * as a prefix, it becomes a partial match.

  $ parade keys *Test*

  If you do not specify any queries, display all keys.
`
	queryExampleGet = `  get command supports exact match, forward match, and partial match.
  parade usually searches for exact matches, and shows only the value.

  $ parade keys /MyService/Test
  value

  Use * as a postfix, the search will be done as a forward match and shows matched keys and values.

  $ parade keys /MyService*
  /MyService/Test  value

  Furthermore, also use * as a prefix, it becomes a partial match.

  $ parade keys *Test*
  /MyService/Test   value
  /MyService2/Test  value2
`
)

type GlobalOption struct {
	Profile   string
	Region    string
	IsNoColor bool

	SSMManager ssmctl.SSMManager

	Out    io.Writer
	ErrOut io.Writer
}

func NewRootCommand(outWriter, errWriter io.Writer) (*cobra.Command, error) {
	o := &GlobalOption{}
	cmd := &cobra.Command{
		Use:     "parade",
		Short:   "parade is a simple AWS SSM parameters CLI",
		Version: VERSION,
		Long:    longDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Use subcommand: keys, get, set, del")
		},
	}
	cmd.PersistentFlags().StringVarP(&o.Profile, "profile", "p", "", "AWS profile")
	// aws-sdk-go does not support the AWS_DEFAULT_REGION environment variable
	cmd.PersistentFlags().StringVar(&o.Region, "region", os.Getenv("AWS_DEFAULT_REGION"), "AWS region")
	cmd.PersistentFlags().BoolVar(&o.IsNoColor, "no-color", false, "Turn off colored output")

	o.Out = outWriter
	o.ErrOut = errWriter
	cmd.SetOut(outWriter)
	cmd.SetErr(errWriter)

	cmd.AddCommand(
		newKeysCommand(o),
		newGetCommand(o),
		newSetCommand(o),
		newDelCommand(o),
	)

	return cmd, nil
}

// Execute executes the root command
func Execute() error {
	o := os.Stdout
	e := os.Stderr

	rootCmd, err := NewRootCommand(o, e)
	if err != nil {
		return err
	}

	return rootCmd.Execute()
}

func initializeCredential(profile, region string) (ssmctl.SSMManager, error) {
	ssmManager, err := ssmctl.New(profile, region)
	if err != nil {
		return nil, fmt.Errorf(ErrMsgAWSProfileNotValid)
	}
	return ssmManager, nil
}
