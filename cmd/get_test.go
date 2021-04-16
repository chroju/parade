package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/chroju/parade/ssmctl"
)

const (
	usageGetHelp = "  -h, --help      help for get\n\n"
)

func Test_getCommand(t *testing.T) {
	voidCmd := newGetCommand(
		&GlobalOption{
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	)
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
			wantOutWriter: "dev_value1\n",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "one arg for Equals with encrypted",
			command:       "/service1/dev/key2",
			wantOutWriter: "(encrypted)\n",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "one arg for Equals with decrypted",
			command:       "/service1/dev/key2 -d",
			wantOutWriter: "dev_value2\n",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "one arg for BeginsWith",
			command:       "service1/*",
			wantOutWriter: "/service1/dev/key1   dev_value1\n/service1/dev/key2   (encrypted)\n/service1/prod/key1  prod_value1\n/service1/prod/key2  (encrypted)\n",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "one arg for Contains",
			command:       "*prod*",
			wantOutWriter: "/service1/prod/key1  prod_value1\n/service1/prod/key2  (encrypted)\n",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "one arg for not found",
			command:       "not_found",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "two args",
			command:       "dev prod",
			wantOutWriter: fmt.Sprintf("Error: accepts 1 arg(s), received 2\n%s%s", voidCmd.UsageString(), usageGetHelp),
			wantErrWriter: "",
			wantErr:       true,
		},
		{
			name:          "no args",
			command:       "",
			wantOutWriter: fmt.Sprintf("Error: accepts 1 arg(s), received 0\n%s%s", voidCmd.UsageString(), usageGetHelp),
			wantErrWriter: "",
			wantErr:       true,
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

			cmd := newGetCommand(o)
			if tt.command != "" {
				cmd.SetArgs(strings.Split(tt.command, " "))
			}

			if err := cmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutWriter := outWriter.String(); gotOutWriter != tt.wantOutWriter {
				t.Errorf("get() = %v, want %v", gotOutWriter, tt.wantOutWriter)
			}
			if gotErrWriter := errWriter.String(); gotErrWriter != tt.wantErrWriter {
				t.Errorf("get() = %v, want %v", gotErrWriter, tt.wantErrWriter)
			}
		})
	}
}
