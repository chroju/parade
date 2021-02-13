package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/chroju/parade/ssmctl"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// VERSION is cli tool version
const VERSION = "0.2.0"

const (
	ErrMsgAWSProfileNotValid = `ERROR: AWS credential is not valid.

Use the --profile and --region options, or set the access keys and region in the environment variables.`
)

var (
	// StdWriter is the io.Writer for standard output
	StdWriter io.Writer
	// ErrWriter is the io.Writer for error output
	ErrWriter io.Writer

	profile    string
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

	var err error
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "AWS profile")
	ssmManager, err = ssmctl.New(profile)
	if err != nil {
		fmt.Fprintln(ErrWriter, color.RedString(ErrMsgAWSProfileNotValid))
		return err
	}

	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(KeysCommand)
	rootCmd.AddCommand(GetCommand)
	rootCmd.AddCommand(SetCommand)
	rootCmd.AddCommand(DelCommand)
}
