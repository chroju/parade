package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/chroju/parade/ssmctl"
)

const (
	usageSetHelp = "  -h, --help      help for set\n\n"
)

func Test_setCommand(t *testing.T) {
	tests := []struct {
		name          string
		command       string
		wantOutWriter string
		wantErrWriter string
		wantErr       bool
	}{
		{
			name:          "two args",
			command:       "/service1/dev/key3 value3",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "two args not overwrite",
			command:       "/service1/dev/key1 value1",
			wantOutWriter: "",
			wantErrWriter: "WARN: `/service1/dev/key1` already exists.\nOverwrite `/service1/dev/key1` (value: dev_value1) ? (Y/n)\n",
			wantErr:       false,
		},
		{
			name:          "two args overwrite",
			command:       "/service1/dev/key1 value1 -f",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "two args encryption",
			command:       "/service1/dev/key3 value3 --encrypt",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "two args overwrite and encryption",
			command:       "/service1/dev/key1 value1 -f --encrypt",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "two args encryption and specify KMS key ID",
			command:       "/service1/dev/key3 value3 --encrypt --kms-key-id key",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "two args specify KMS key ID without encryption",
			command:       "/service1/dev/key3 value3 --kms-key-id key",
			wantOutWriter: "",
			wantErrWriter: "Failed to execute PutParameter API.\nKMS Key ID must be used with SecureString type.",
			wantErr:       true,
		},
		{
			name:          "one arg",
			command:       "/service1/dev/key1",
			wantOutWriter: "",
			wantErrWriter: "accepts 2 arg(s), received 1",
			wantErr:       true,
		},
		{
			name:          "no args",
			command:       "",
			wantOutWriter: "",
			wantErrWriter: "accepts 2 arg(s), received 0",
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

			cmd := newSetCommand(o)
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
