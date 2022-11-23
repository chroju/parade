package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/chroju/parade/ssmctl"
)

const (
	usageKeysHelp = "  -h, --help      help for keys\n\n"
	// TODO
	usageKeys = `Usage:
  keys [query] [flags]

Examples:
  keys command supports exact match, forward match, and partial match.
  It usually searches for exact matches.

  $ parade keys /MyService/Test

  Use * as a postfix, the search will be done as a forward match.

  $ parade keys /MyService*

  Furthermore, also use * as a prefix, it becomes a partial match.

  $ parade keys *Test*

  If you do not specify any queries, display all keys.


Flags:
  -h, --help       help for keys
      --no-types   Turn off parameter type shows`
)

func Test_keysCommand(t *testing.T) {
	type args struct {
		args       []string
		ssmManager ssmctl.SSMManager
	}
	tests := []struct {
		name          string
		command       string
		wantOutWriter string
		wantErrWriter string
		wantErr       bool
	}{
		{
			name:          "one arg for Equals",
			command:       "/service1/dev/key1",
			wantOutWriter: "/service1/dev/key1  Type: String\n",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "one arg for BeginsWith",
			command:       "/service1/*",
			wantOutWriter: "/service1/dev/key1   Type: String\n/service1/dev/key2   Type: SecureString\n/service1/prod/key1  Type: String\n/service1/prod/key2  Type: SecureString\n",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "one arg for Contains",
			command:       "*dev*",
			wantOutWriter: "/service1/dev/key1  Type: String\n/service1/dev/key2  Type: SecureString\n",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "one arg for no match",
			command:       "no_match",
			wantOutWriter: "",
			wantErrWriter: fmt.Sprintf("%s\nParameterNotFound", ErrMsgDescribeParameters),
			wantErr:       true,
		},
		{
			name:          "no args",
			command:       "",
			wantOutWriter: "/service1/dev/key1   Type: String\n/service1/dev/key2   Type: SecureString\n/service1/prod/key1  Type: String\n/service1/prod/key2  Type: SecureString\n",
			wantErrWriter: "",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ssmManager := ssmctl.NewMockSSMManager()
			outWriter := &bytes.Buffer{}
			errWriter := &bytes.Buffer{}

			o := &GlobalOption{
				SSMManager: ssmManager,
				Out:        outWriter,
				ErrOut:     errWriter,
			}

			cmd := newKeysCommand(o)
			if tt.command != "" {
				cmd.SetArgs(strings.Split(tt.command, " "))
			}

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if err.Error() != tt.wantErrWriter {
					t.Errorf("get() = %v, want %v", err.Error(), tt.wantErrWriter)
				}
			} else {
				if gotOutWriter := outWriter.String(); gotOutWriter != tt.wantOutWriter {
					t.Errorf("get() = %v, want %v", gotOutWriter, tt.wantOutWriter)
				}
			}
		})
	}
}
